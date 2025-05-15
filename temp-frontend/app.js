const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const infoBox = document.getElementById('info-box');
const adminApiUrl = 'https://utmn-map.zetoqqq.ru/adminapi';
const publicApiUrl = 'https://utmn-map.zetoqqq.ru/publicapi';

let allData = null;
let currentFloor = null;
let offsetX = 0, offsetY = 0;
let drag = false, lastX = 0, lastY = 0;
let scale = 1;
let visibleObjects = [];
let selectedObject = null;
let currentBuildingId = null;

async function init() {
    try {
        const res = await fetch(`${publicApiUrl}/api/buildings`);
        const data = await res.json();
        if (!Array.isArray(data.buildings) || data.buildings.length === 0) {
            throw new Error("Нет доступных зданий");
        }
        currentBuildingId = data.buildings[0].id; // выбираем первое здание
        console.log("Выбранное здание:", currentBuildingId);
        loadBuildingObjects(currentBuildingId);
    } catch (err) {
        console.error("Ошибка при загрузке зданий:", err);
        infoBox.innerHTML = "Не удалось загрузить список зданий.";
    }
}

function loadBuildingObjects(buildingId) {
    fetch(`${publicApiUrl}/api/buildings/${buildingId}/objects`)
        .then(res => res.json())
        .then(data => {
            const result = data.objects; // извлекаем объект с building, floors, background
            console.log("Извлечённые данные:", result);
            if (!Array.isArray(result.floors)) throw new Error("Ожидался массив этажей");
            allData = result;
            createFloorButtons(allData);
            resizeCanvas();
        })
        .catch(err => {
            console.error("Ошибка при загрузке объектов:", err);
            infoBox.innerHTML = "Ошибка загрузки данных с сервера.";
        });
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

// Обновленный обработчик клика по canvas
canvas.addEventListener('click', e => {
    const rect = canvas.getBoundingClientRect();
    const clickX = (e.clientX - rect.left - offsetX) / scale;
    const clickY = (e.clientY - rect.top - offsetY) / scale;

    // Сначала ищем клик по двери
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

    // Если клик не по двери — ищем клик по объекту
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
    // data.floors — это массив этажей, каждый имеет поле floor.name
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

// Функция для корректировки яркости цвета (используется для объектов и дверей)
function adjustColor(color, factor) {
    // Ожидается формат rgba(r, g, b, a)
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

/**
 * Проверяет, попадает ли точка (px,py) в повернутый прямоугольник двери.
 * Переводим точку в локальные координаты двери (относительно её оси поворота),
 * затем проверяем, лежит ли она в пределах [0, door.width] и [0, door.height].
 */
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

/**
 * Анимация (открытия/закрытия) двери.
 * Плавно поворачивает дверь от текущего угла до targetAngle за duration миллисекунд.
 */
function animateDoor(door, targetAngle) {
    const duration = 500; // длительность анимации в мс
    const startAngle = door.angle || 0;
    const startTime = performance.now();

    function step() {
        const now = performance.now();
        const elapsed = now - startTime;
        const progress = Math.min(elapsed / duration, 1);
        door.angle = startAngle + (targetAngle - startAngle) * progress;
        visualize(allData); // перерисовываем сцену
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

    // Определяем, какие этажи отрисовывать: если выбран конкретный этаж – только его, иначе все
    const floorData = currentFloor
        ? allData.floors.find(f => f.floor.name === currentFloor)
        : null;
    const floorsToRender = floorData ? [floorData] : allData.floors;

    for (const floor of floorsToRender) {
        // Рисуем фон этажа (background)
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

            // Для типа "cabinet" используем непрозрачный цвет
            let color = {
                'cabinet': 'rgba(0, 128, 255, 1)', // без прозрачности
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

    ctx.restore();
}

// Обработчик клика по canvas
canvas.addEventListener('click', e => {
    const rect = canvas.getBoundingClientRect();
    const clickX = (e.clientX - rect.left - offsetX) / scale;
    const clickY = (e.clientY - rect.top - offsetY) / scale;

    // Проверяем, был ли клик по двери (и запоминаем родительский объект)
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

    // Если клик не по двери, проверяем попадание по объектам
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

// 🔥 Запуск
init();
