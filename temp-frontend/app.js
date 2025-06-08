const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const infoBox = document.getElementById('info-box');
const adminApiUrl = 'https://utmn-map.zetoqqq.ru/adminapi';
const publicApiUrl = 'https://utmn-map.zetoqqq.ru/publicapi';

// 1. Bearer-токен для adminApi
const ADMIN_API_TOKEN = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDkzNzU5NzUsImlhdCI6MTc0OTM3NDE3NSwianRpIjoiZTRlOTM0YmMtZTM5ZS00MDYyLThhNzEtMDFkNGU3N2U3NjE2Iiwicm9sZXMiOlsidXNlciJdLCJzdWIiOiI5Yjg0Y2Q3Yy01N2FkLTQzY2UtOTljNi02MTViODdjZTdiY2QifQ.y3l7oW6fQJTYa0unE8etsFMMS3_ls-Lw_f2_K0NxqdY';

let allData = null;
let currentFloor = null;
let offsetX = 0, offsetY = 0;
let drag = false, lastX = 0, lastY = 0;
let scale = 1;
let visibleObjects = [];
let selectedObject = null;
let currentBuildingId = null;

// Новые глобальные переменные для графа:
let graphNodes = [];        // объекты вида { floor_id, id, type, x, y }
let graphConnections = [];  // объекты вида { from_id, to_id, weight }

// Новая переменная для хранения выбранных узлов (максимум 2)
let selectedNodes = [];

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
        updatedFields[key] = el.value;
    });

    fetch(`${adminApiUrl}/api/objects/${selectedObject.id}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': ADMIN_API_TOKEN
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
            infoBox.innerHTML = `<span style="color:red;">Ошибка при сохранении: ${err.message}</span>`;
        });
}

function resizeCanvas() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    if (allData) visualize(allData);
}
window.addEventListener('resize', resizeCanvas);

// Перетаскивание и масштабирование
canvas.addEventListener('mousedown', e => {
    drag = true;
    lastX = e.clientX;
    lastY = e.clientY;
});
canvas.addEventListener('mouseup', () => drag = false);
canvas.addEventListener('mouseleave', () => drag = false);
canvas.addEventListener('mousemove', e => {
    if (!drag) return;
    offsetX += e.clientX - lastX;
    offsetY += e.clientY - lastY;
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
    scale *= e.deltaY < 0 ? zoomFactor : 1 / zoomFactor;
    offsetX = x - ((x - offsetX) / prevScale) * scale;
    offsetY = y - ((y - offsetY) / prevScale) * scale;
    visualize(allData);
}, { passive: false });

// Обработка клика: двери, объекты и узлы графа
canvas.addEventListener('click', e => {
    const rect = canvas.getBoundingClientRect();
    const clickX = (e.clientX - rect.left - offsetX) / scale;
    const clickY = (e.clientY - rect.top - offsetY) / scale;

    // Сначала — двери (старый код)
    for (const object of visibleObjects) {
        if (Array.isArray(object.doors)) {
            for (const door of object.doors) {
                if (isPointInRotatedRect(clickX, clickY, door)) {
                    showDoorInfo(door, object);
                    return;
                }
            }
        }
    }

    // Затем — объекты
    for (const obj of visibleObjects) {
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

    // Наконец — узлы графа
    const node = findNodeUnderCursor(clickX, clickY);
    if (node) {
        // добавляем в selectedNodes (макс. 2, без дубликатов)
        if (!selectedNodes.find(n => n.id === node.id)) {
            if (selectedNodes.length === 2) selectedNodes.shift();
            selectedNodes.push(node);
        }
        showNodeInfo(node);
        return;
    }

    infoBox.innerHTML = 'Ничего не найдено под курсором.';
});

// Редактирование обычного объекта
async function showObjectInfo(obj) {
    selectedObject = obj;
    selectedNodes = []; // сброс выбора узлов
    let html = `<b>Редактирование объекта:</b><br>`;
    html += `<b>id:</b> <div style="font-family: monospace;">${obj.id}</div>`;
    html += `<b>coords:</b> x=${obj.x}; y=${obj.y}; w=${obj.width}; h=${obj.height}<br>`;
    html += `<b>name:</b> <input data-key="name" value="${obj.name || ''}" style="width:100%"><br>`;
    html += `<b>alias:</b> <input data-key="alias" value="${obj.alias || ''}" style="width:100%"><br>`;
    html += `<b>description:</b> <input data-key="description" value="${obj.description || ''}" style="width:100%"><br>`;
    // Тип объекта
    let opts = `<option value="">-- выберите тип --</option>`;
    try {
        const res = await fetch(`${publicApiUrl}/api/categories`);
        const cats = (await res.json()).categories;
        opts += cats.map(c => `<option value="${c}" ${c === obj.object_type ? 'selected' : ''}>${c}</option>`).join('');
    } catch {
        opts += `<option disabled>Ошибка загрузки</option>`;
    }
    html += `<b>object_type:</b> <select data-key="object_type" style="width:100%">${opts}</select><br>`;
    html += `<button onclick="saveObject()">Сохранить</button>`;
    infoBox.innerHTML = html;
}

// Инфо о двери
function showDoorInfo(door, parent) {
    selectedNodes = []; // сброс выбора узлов
    let html = `<b>Информация о двери:</b><br>`;
    html += `door id: <span style="font-family: monospace;">${door.id}</span><br>`;
    html += `x=${door.x}, y=${door.y}, w=${door.width}, h=${door.height}<br>`;
    html += `angle=${door.angle || 0}<hr>`;
    html += `<b>Родительский объект:</b><br>`;
    html += `object id: <span style="font-family: monospace;">${parent.id}</span><br>`;
    html += `name: ${parent.name || '?'}, alias: ${parent.alias || ''}`;
    infoBox.innerHTML = html;
}

// Инфо о выбранном узле и кнопка «Создать связь»
function showNodeInfo(node) {
    selectedObject = null;
    let html = `<b>Узел графа:</b><br>`;
    html += `id: <span style="font-family: monospace;">${node.id}</span><br>`;
    html += `type: ${node.type}<br>`;
    if (selectedNodes.length === 2) {
        html += `<button id="create-connection-btn">Создать связь</button>`;
    }
    infoBox.innerHTML = html;
    const btn = document.getElementById('create-connection-btn');
    if (btn) btn.onclick = createConnection;
}

// Отправка POST-запроса для создания связи
async function createConnection() {
    const [fromNode, toNode] = selectedNodes;
    const payload = { from_id: fromNode.id, to_id: toNode.id, weight: 1 };
    try {
        const res = await fetch(
            `${adminApiUrl}/api/buildings/${currentBuildingId}/route/connections`,
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': ADMIN_API_TOKEN
                },
                body: JSON.stringify(payload)
            }
        );
        if (!res.ok) throw new Error(`Status ${res.status}`);
        const { connection } = await res.json();
        graphConnections.push({
            from_id: connection.from_id,
            to_id: connection.to_id,
            weight: connection.weight
        });
        visualize(allData);
        selectedNodes = [];
        infoBox.innerHTML = `<span style="color:green;">Связь создана!</span>`;
    } catch (err) {
        console.error(err);
        infoBox.innerHTML = `<span style="color:red;">Ошибка: ${err.message}</span>`;
    }
}

// Функция поиска узла под курсором (радиус 5px)
function findNodeUnderCursor(x, y) {
    const floorData = currentFloor
        ? allData.floors.find(f => f.floor.name === currentFloor)
        : null;
    const floorId = floorData ? floorData.floor.id : null;
    return graphNodes
        .filter(n => n.floor_id === floorId)
        .find(n => {
            const dx = x - n.x;
            const dy = y - n.y;
            return dx * dx + dy * dy <= 5 * 5;
        }) || null;
}

// --- Остальной код без изменений ---

function createFloorButtons(data) {
    const floors = data.floors.map(f => f.floor.name);
    const container = document.getElementById('floor-buttons');
    container.innerHTML = '';
    const allBtn = document.createElement('button');
    allBtn.textContent = 'Все этажи';
    allBtn.onclick = () => { currentFloor = null; visualize(allData); };
    container.appendChild(allBtn);
    floors.forEach(floor => {
        const btn = document.createElement('button');
        btn.textContent = floor;
        btn.onclick = () => { currentFloor = floor; visualize(allData); };
        container.appendChild(btn);
    });
}

function adjustColor(color, factor) {
    const parts = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+),?\s*([\d.]+)?\)/);
    if (parts) {
        let [r, g, b, a] = parts.slice(1).map((v, i) => i < 3 ? Math.round(v * factor) : parseFloat(v || 1));
        return `rgba(${r},${g},${b},${a})`;
    }
    return color;
}
function lightenColor(c, f = 1.2) { return adjustColor(c, f); }
function darkenColor(c, f = 0.8) { return adjustColor(c, f); }

function isPointInRotatedRect(px, py, door) {
    const a = door.angle || 0;
    const dx = px - door.x, dy = py - door.y;
    const cos = Math.cos(-a), sin = Math.sin(-a);
    const lx = dx * cos - dy * sin, ly = dx * sin + dy * cos;
    return lx >= 0 && lx <= door.width && ly >= 0 && ly <= door.height;
}

function animateDoor(door, targetAngle) {
    const start = performance.now();
    const from = door.angle || 0;
    (function step(now) {
        const t = Math.min((now - start) / 500, 1);
        door.angle = from + (targetAngle - from) * t;
        visualize(allData);
        if (t < 1) requestAnimationFrame(step);
    })(start);
}

function visualize(data) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.save();
    ctx.translate(offsetX, offsetY);
    ctx.scale(scale, scale);
    visibleObjects = [];

    const floorData = currentFloor
        ? allData.floors.find(f => f.floor.name === currentFloor)
        : null;
    const floors = floorData ? [floorData] : allData.floors;

    for (const fl of floors) {
        fl.background.forEach(bg => {
            const pts = bg.points.sort((a, b) => a.order - b.order);
            if (pts.length > 1) {
                ctx.beginPath();
                ctx.moveTo(pts[0].x, pts[0].y);
                pts.slice(1).forEach(p => ctx.lineTo(p.x, p.y));
                ctx.closePath();
                ctx.fillStyle = '#f0f0f0'; ctx.fill();
                ctx.strokeStyle = '#999'; ctx.stroke();
            }
        });
        fl.objects.forEach(obj => {
            const { x, y, width, height, object_type, name, doors } = obj;
            const colors = {
                'cabinet': 'rgba(0,128,255,1)',
                'wardrobe': 'rgba(255,165,0,0.5)',
                'woman-toilet': 'rgba(255,192,203,0.5)',
                'man-toilet': 'rgba(144,238,144,0.5)',
                'gym': 'rgba(128,0,128,0.5)',
            };
            ctx.fillStyle = colors[object_type] || 'rgba(200,200,200,0.5)';
            ctx.fillRect(x, y, width, height);
            ctx.strokeStyle = 'black'; ctx.strokeRect(x, y, width, height);
            ctx.fillStyle = 'black'; ctx.font = '14px Arial';
            ctx.fillText(name || '???', x + 5, y + 15);
            if (Array.isArray(doors)) {
                doors.forEach(d => {
                    ctx.fillStyle = 'red';
                    ctx.fillRect(d.x, d.y, d.width, d.height);
                    ctx.strokeRect(d.x, d.y, d.width, d.height);
                });
            }
            visibleObjects.push(obj);
        });
    }

    // Рисуем граф на текущем этаже
    if (currentFloor && floorData) {
        const fid = floorData.floor.id;
        const nodesOnF = graphNodes.filter(n => n.floor_id === fid);
        ctx.strokeStyle = 'green'; ctx.lineWidth = 2;
        graphConnections.forEach(c => {
            const f = nodesOnF.find(n => n.id === c.from_id);
            const t = nodesOnF.find(n => n.id === c.to_id);
            if (f && t) {
                ctx.beginPath();
                ctx.moveTo(f.x, f.y);
                ctx.lineTo(t.x, t.y);
                ctx.stroke();
            }
        });
        ctx.fillStyle = 'red';
        nodesOnF.forEach(n => {
            ctx.beginPath();
            ctx.arc(n.x, n.y, 2, 0, Math.PI * 2);
            ctx.fill();
        });
    }

    ctx.restore();
}

init();
