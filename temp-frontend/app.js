const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const infoBox = document.getElementById('info-box');
const adminApiUrl = 'http://localhost:8001';
const publicApiUrl = 'http://localhost:8000';

let allData = null;
let currentFloor = null;
let offsetX = 0, offsetY = 0;
let drag = false, lastX = 0, lastY = 0;
let scale = 1;
let visibleObjects = [];
let selectedObject = null;
let currentBuildingId = null;

// Новые глобальные переменные для графа:
let graphNodes = [];       // объекты вида { floor_id, id, type, x, y }
let graphConnections = []; // объекты вида { from_id, to_id, weight }

async function init() {
    try {
        const res = await fetch(`${publicApiUrl}/api/buildings`);
        const data = await res.json();
        if (!Array.isArray(data.buildings) || data.buildings.length === 0) {
            throw new Error("Нет доступных зданий");
        }
        currentBuildingId = data.buildings[0].id;
        console.log("Выбранное здание:", currentBuildingId);

        await loadBuildingObjects(currentBuildingId);
        await loadGraphData(currentBuildingId);
        resizeCanvas();
    } catch (err) {
        console.error("Ошибка при инициализации:", err);
        infoBox.innerHTML = "Не удалось загрузить данные.";
    }
}

function loadBuildingObjects(buildingId) {
    return fetch(`${publicApiUrl}/api/buildings/${buildingId}/objects`)
        .then(res => res.json())
        .then(data => {
            const result = data.objects;
            console.log("Извлечённые данные по объектам:", result);
            if (!Array.isArray(result.floors)) throw new Error("Ожидался массив этажей");
            allData = result;
            createFloorButtons(allData);
        })
        .catch(err => {
            console.error("Ошибка при загрузке объектов:", err);
            infoBox.innerHTML = "Ошибка загрузки данных с сервера.";
        });
}

function loadGraphData(buildingId) {
    const nodesPromise = fetch(`${publicApiUrl}/api/buildings/${buildingId}/graph/nodes`)
        .then(res => res.json())
        .then(data => {
            if (Array.isArray(data.nodes)) {
                graphNodes = data.nodes;
            } else {
                console.warn("Ожидался массив nodes");
            }
        })
        .catch(err => {
            console.error("Ошибка при загрузке узлов (nodes):", err);
        });

    const connsPromise = fetch(`${publicApiUrl}/api/buildings/${buildingId}/connections`)
        .then(res => res.json())
        .then(data => {
            if (Array.isArray(data.connections)) {
                graphConnections = data.connections;
            } else {
                console.warn("Ожидался массив connections");
            }
        })
        .catch(err => {
            console.error("Ошибка при загрузке связей (connections):", err);
        });

    return Promise.all([nodesPromise, connsPromise]);
}

function saveObject() {
    if (!selectedObject || !selectedObject.id) return;

    const inputs = infoBox.querySelectorAll('[data-key]');
    const allowedKeys = ['name', 'object_type', 'alias', 'description'];
    const updatedFields = {};

    inputs.forEach(el => {
        const key = el.dataset.key;
        if (!allowedKeys.includes(key)) return;
        const value = el.value;
        updatedFields[key] = value;
    });

    fetch(`${adminApiUrl}/api/objects/${selectedObject.id}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedFields)
    })
        .then(res => {
            if (!res.ok) throw new Error('Ошибка при сохранении');
            return res.json();
        })
        .then(updated => {
            Object.assign(selectedObject, updated);
            visualize(allData);
        })
        .catch(err => {
            console.error(err);
        });
}

function resizeCanvas() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    if (allData) visualize(allData);
}
window.addEventListener('resize', resizeCanvas);

canvas.addEventListener('mousedown', e => {
    drag = true;
    lastX = e.clientX;
    lastY = e.clientY;
});
canvas.addEventListener('mouseup', () => drag = false);
canvas.addEventListener('mouseleave', () => drag = false);
canvas.addEventListener('mousemove', e => {
    if (!drag) return;
    const dx = e.clientX - lastX;
    const dy = e.clientY - lastY;
    offsetX += dx;
    offsetY += dy;
    lastX = e.clientX;
    lastY = e.clientY;
    visualize(allData);
});
canvas.addEventListener('wheel', e => {
    e.preventDefault();
    const zoomFactor = 1.05;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const prevScale = scale;
    if (e.deltaY < 0) {
        scale *= zoomFactor;
    } else {
        scale /= zoomFactor;
    }
    offsetX = x - ((x - offsetX) / prevScale) * scale;
    offsetY = y - ((y - offsetY) / prevScale) * scale;
    visualize(allData);
}, { passive: false });

canvas.addEventListener('click', e => {
    const rect = canvas.getBoundingClientRect();
    const clickX = (e.clientX - rect.left - offsetX) / scale;
    const clickY = (e.clientY - rect.top - offsetY) / scale;

    let doorClickedInfo = null;
    for (const object of visibleObjects) {
        if (Array.isArray(object.doors)) {
            for (const door of object.doors) {
                if (isPointInRotatedRect(clickX, clickY, door)) {
                    doorClickedInfo = { door, parent: object };
                    break;
                }
            }
        }
        if (doorClickedInfo) break;
    }
    if (doorClickedInfo) {
        showDoorInfo(doorClickedInfo.door, doorClickedInfo.parent);
        return;
    }

    for (let obj of visibleObjects) {
        if (
            clickX >= obj.x &&
            clickX <= obj.x + obj.width &&
            clickY >= obj.y &&
            clickY <= obj.y + obj.height
        ) {
            showObjectInfo(obj);
            return;
        }
    }
    infoBox.innerHTML = 'Объект не найден. Нажмите на элемент на карте.';
});

async function showObjectInfo(obj) {
    selectedObject = obj;

    let html = `<b>Редактирование объекта:</b><br>`;
    html += `<b>id:</b> <div style="font-family: monospace; user-select: all">${obj.id}</div>`;
    html += `<b>coords:</b> <div style="font-family: monospace; user-select: all">x=${obj.x}; y=${obj.y}<br>w=${obj.width}; h=${obj.height}</div>`;
    html += `<b>name:</b> <input data-key="name" value="${obj.name || ''}" style="width: 100%;"><br>`;
    html += `<b>alias:</b> <input data-key="alias" value="${obj.alias || ''}" style="width: 100%;"><br>`;
    html += `<b>description:</b> <input data-key="description" value="${obj.description || ''}" style="width: 100%;"><br>`;

    let optionsHTML = `<option value="">-- выберите тип --</option>`;
    try {
        const res = await fetch(`${publicApiUrl}/api/categories`);
        const data = await res.json();
        const categories = data.categories;
        optionsHTML += categories.map(cat =>
            `<option value="${cat}" ${cat === obj.object_type ? 'selected' : ''}>${cat}</option>`
        ).join('');
    } catch (err) {
        console.warn('Не удалось загрузить категории типов:', err);
        optionsHTML += `<option disabled>Ошибка загрузки типов</option>`;
    }

    html += `<b>object_type:</b> 
        <select data-key="object_type" style="width: 100%;">${optionsHTML}</select><br>`;
    html += `<button onclick="saveObject()">Сохранить</button>`;

    infoBox.innerHTML = html;
}

function showDoorInfo(door, parent) {
    let html = `<b>Информация о двери:</b><br>`;
    html += `<b>door id:</b> <div style="font-family: monospace; user-select: all">${door.id}</div>`;
    html += `<b>Координаты:</b> <div style="font-family: monospace;">x=${door.x}, y=${door.y}, w=${door.width}, h=${door.height}</div>`;
    html += `<b>Угол (angle):</b> <div style="font-family: monospace;">${door.angle || 0}</div>`;

    html += `<hr><b>Объект:</b><br>`;
    html += `<b>object id:</b> <div style="font-family: monospace;">${parent.id}</div>`;
    html += `<b>name:</b> <div>${parent.name || '???'}</div>`;
    html += `<b>alias:</b> <div>${parent.alias || ''}</div>`;

    infoBox.innerHTML = html;
}

function createFloorButtons(data) {
    const floors = data.floors.map(f => f.floor.name);
    const container = document.getElementById('floor-buttons');
    container.innerHTML = '';

    const allBtn = document.createElement('button');
    allBtn.textContent = 'Все этажи';
    allBtn.onclick = () => {
        currentFloor = null;
        visualize(allData);
    };
    container.appendChild(allBtn);

    floors.forEach(floor => {
        const btn = document.createElement('button');
        btn.textContent = floor;
        btn.onclick = () => {
            currentFloor = floor;
            visualize(allData);
        };
        container.appendChild(btn);
    });
}

function adjustColor(color, factor) {
    const parts = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+),?\s*([\d.]+)?\)/);
    if (parts) {
        let r = Math.min(255, Math.max(0, Math.round(parseFloat(parts[1]) * factor)));
        let g = Math.min(255, Math.max(0, Math.round(parseFloat(parts[2]) * factor)));
        let b = Math.min(255, Math.max(0, Math.round(parseFloat(parts[3]) * factor)));
        let a = parts[4] ? parseFloat(parts[4]) : 1;
        return `rgba(${r}, ${g}, ${b}, ${a})`;
    }
    return color;
}

function lightenColor(color, factor = 1.2) {
    return adjustColor(color, factor);
}

function darkenColor(color, factor = 0.8) {
    return adjustColor(color, factor);
}

function isPointInRotatedRect(px, py, door) {
    const angle = door.angle || 0;
    const dx = px - door.x;
    const dy = py - door.y;
    const cos = Math.cos(-angle);
    const sin = Math.sin(-angle);
    const localX = dx * cos - dy * sin;
    const localY = dx * sin + dy * cos;
    return localX >= 0 && localX <= door.width && localY >= 0 && localY <= door.height;
}

function animateDoor(door, targetAngle) {
    const duration = 500;
    const startAngle = door.angle || 0;
    const startTime = performance.now();

    function step() {
        const now = performance.now();
        const elapsed = now - startTime;
        const progress = Math.min(elapsed / duration, 1);
        door.angle = startAngle + (targetAngle - startAngle) * progress;
        visualize(allData);
        if (progress < 1) {
            requestAnimationFrame(step);
        }
    }
    requestAnimationFrame(step);
}

function visualize(data) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.save();
    ctx.translate(offsetX, offsetY);
    ctx.scale(scale, scale);

    visibleObjects = [];

    // Определяем, какие этажи отрисовывать
    const floorData = currentFloor
        ? allData.floors.find(f => f.floor.name === currentFloor)
        : null;
    const floorsToRender = floorData ? [floorData] : allData.floors;

    for (const floor of floorsToRender) {
        // Рисуем фон этажа
        for (const bg of floor.background) {
            const sortedPoints = bg.points.sort((a, b) => a.order - b.order);
            if (sortedPoints.length > 1) {
                ctx.beginPath();
                ctx.moveTo(sortedPoints[0].x, sortedPoints[0].y);
                for (let i = 1; i < sortedPoints.length; i++) {
                    ctx.lineTo(sortedPoints[i].x, sortedPoints[i].y);
                }
                ctx.closePath();
                ctx.fillStyle = '#f0f0f0';
                ctx.fill();
                ctx.strokeStyle = '#999';
                ctx.stroke();
            }
        }

        // Рисуем объекты этажа
        for (const object of floor.objects) {
            const { x, y, width, height } = object;
            let color = {
                'cabinet': 'rgba(0, 128, 255, 1)',
                'wardrobe': 'rgba(255, 165, 0, 0.5)',
                'woman-toilet': 'rgba(255, 192, 203, 0.5)',
                'man-toilet': 'rgba(144, 238, 144, 0.5)',
                'gym': 'rgba(128, 0, 128, 0.5)',
            }[object.object_type] || 'rgba(200, 200, 200, 0.5)';

            ctx.fillStyle = color;
            ctx.fillRect(x, y, width, height);
            ctx.strokeStyle = 'black';
            ctx.strokeRect(x, y, width, height);

            ctx.fillStyle = 'black';
            ctx.font = '14px Arial';
            ctx.fillText(object.name || '???', x + 5, y + 15);

            if (Array.isArray(object.doors)) {
                for (const door of object.doors) {
                    ctx.fillStyle = 'red';
                    ctx.fillRect(door.x, door.y, door.width, door.height);
                    ctx.strokeRect(door.x, door.y, door.width, door.height);
                }
            }

            visibleObjects.push(object);
        }
    }

    // -------------- Рисуем граф только для выбранного этажа ----------------

    if (currentFloor && floorData) {
        // Найдём ID выбранного этажа
        const selectedFloorId = floorData.floor.id;

        // Фильтруем узлы по floor_id
        const nodesOnFloor = graphNodes.filter(n => n.floor_id === selectedFloorId);

        // Рисуем связи только между узлами этого же этажа
        ctx.strokeStyle = 'green';
        ctx.lineWidth = 2;
        graphConnections.forEach(conn => {
            const fromNode = nodesOnFloor.find(n => n.id === conn.from_id);
            const toNode = nodesOnFloor.find(n => n.id === conn.to_id);
            if (fromNode && toNode) {
                ctx.beginPath();
                ctx.moveTo(fromNode.x, fromNode.y);
                ctx.lineTo(toNode.x, toNode.y);
                ctx.stroke();
            }
        });

        // Рисуем узлы (красные кружки радиусом 5px) только на этом этаже
        ctx.fillStyle = 'red';
        nodesOnFloor.forEach(node => {
            ctx.beginPath();
            ctx.arc(node.x, node.y, 5, 0, Math.PI * 2);
            ctx.fill();
        });
    }

    ctx.restore();
}

// Запуск
init();
