const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const infoBox = document.getElementById('info-box');

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
        const res = await fetch('http://localhost:8000/api/buildings');
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
    fetch(`http://localhost:8000/api/buildings/${buildingId}/objects`)
        .then(res => res.json())
        .then(data => {
            if (!Array.isArray(data.objects)) throw new Error("Ожидался массив объектов");

            allData = data;
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

    fetch(`http://localhost:8000/api/objects/${selectedObject.id}`, {
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
            alert('Объект успешно сохранён!');
        })
        .catch(err => {
            console.error(err);
            alert('Не удалось сохранить объект');
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
    const zoomFactor = 1.1;
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
    html += `<b>name:</b> <input data-key="name" value="${obj.name || ''}" style="width: 100%;"><br>`;
    html += `<b>alias:</b> <input data-key="alias" value="${obj.alias || ''}" style="width: 100%;"><br>`;
    html += `<b>description:</b> <input data-key="description" value="${obj.description || ''}" style="width: 100%;"><br>`;

    let optionsHTML = `<option value="">-- выберите тип --</option>`;
    try {
        const res = await fetch(`http://localhost:8000/api/buildings/${currentBuildingId}/categories`);
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

function createFloorButtons(data) {
    const floors = new Set(data.objects.map(o => o.floor?.name).filter(Boolean));
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

function visualize(data) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.save();
    ctx.translate(offsetX, offsetY);
    ctx.scale(scale, scale);

    visibleObjects = [];

    const objects = currentFloor
        ? data.objects.filter(o => o.floor?.name === currentFloor)
        : data.objects;

    for (const object of objects) {
        const { x, y, width, height } = object;

        let color = {
            'cabinet': 'rgba(0, 128, 255, 0.5)',
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

    ctx.restore();
}

// 🔥 Запуск
init();
