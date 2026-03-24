const CARDS_JSON = "cards.json";
const BASE_URL = 'https://nlr.ru/e-case3/sc2.php/web_gak/gc';
const MAX_ID = 133781;

let CARDS = {};
let currentId = null;
let currentCardNumber = null;

async function loadCardsData() {
    const response = await fetch(CARDS_JSON);
    if (!response.ok) {
        return;
    }
    CARDS = await response.json();
}

function getRandomId() {
    return Math.floor(Math.random() * MAX_ID) + 1;
}

function getRandomCardNumber(id) {
    const maxCN = id in CARDS ? CARDS[id] : 1;
    const randomCN = Math.floor(Math.random() * maxCN) + 1;
    console.log(`For ${id} the max CN is ${maxCN}, showing ${randomCN}`);
    return randomCN
}

function loadCard(id, cardNumber) {
    currentId = id;
    currentCardNumber = cardNumber;
    const img = document.getElementById('card');
    img.src = `${BASE_URL}/${currentId}/${currentCardNumber}`;
    img.alt = `Card ${currentId}/${currentCardNumber}`;
    updateFavouriteButton();
}

function loadRandomCard() {
    const id = getRandomId();
    const cardNumber = getRandomCardNumber(id);
    loadCard(id, cardNumber);
}

function getFavouriteKey(id, cardNumber) {
    return `fav_${id}_${cardNumber}`;
}

function parseFavouriteKey(key) {
    const parts = key.slice(4).split('_');
    const id = parseInt(parts[0]);
    const cardNumber = parseInt(parts[1]);
    return [id, cardNumber];
}

function isFavourited(id, cardNumber) {
    return localStorage.getItem(getFavouriteKey(id, cardNumber)) !== null;
}

function toggleFavourite() {
    const key = getFavouriteKey(currentId, currentCardNumber);

    if (localStorage.getItem(key)) {
        localStorage.removeItem(key);
    } else {
        localStorage.setItem(key, Date.now());
    }

    updateFavouriteButton();
    updateFavouriteCount();
}

function removeFavourite(id, cardNumber) {
    const key = getFavouriteKey(id, cardNumber);
    localStorage.removeItem(key);
}

function updateFavouriteButton() {
    const btn = document.getElementById('favourite');
    btn.textContent = isFavourited(currentId, currentCardNumber) ? '❤️' : '🤍';
}

function getFavouriteKeys() {
    return Array.from({ length: localStorage.length }, (_, i) => localStorage.key(i))
        .filter(key => key?.startsWith('fav_'));
}

function updateFavouriteCount() {
    const count = getFavouriteKeys().length;
    const elLink = document.getElementById('favLink');
    const elCount = document.getElementById('favCount');
    elCount.textContent = `${count}`;
    elLink.style.display = count > 0 ? 'block' : 'none';
    console.log(`favs count updated to ${count}`);
}

function showFavourites() {
    const overlay = document.getElementById('favOverlay');
    const grid = document.getElementById('favGrid');
    grid.innerHTML = '';

    for (const key of getFavouriteKeys()) {
        const [id, cardNumber] = parseFavouriteKey(key);
        const cardDiv = document.createElement('div');
        cardDiv.className = 'favCard';
        cardDiv.onclick = () => {
            hideFavourites();
            loadCard(id, cardNumber);
        };

        const img = document.createElement('img');
        img.src = `${BASE_URL}/${id}/${cardNumber}`;
        img.alt = `Card ${id}/${cardNumber}`;

        const removeBtn = document.createElement('button');
        removeBtn.className = 'favRemove';
        removeBtn.textContent = '×';
        removeBtn.onclick = (e) => {
            e.stopPropagation();
            removeFavourite(id, cardNumber);
            updateFavouriteCount();
            showFavourites();
        };

        cardDiv.appendChild(img);
        cardDiv.appendChild(removeBtn);
        grid.appendChild(cardDiv);
    }

    overlay.style.display = 'flex';
}

function hideFavourites() {
    document.getElementById('favOverlay').style.display = 'none';
}

document.getElementById('newCard').addEventListener('click', loadRandomCard);
document.getElementById('favourite').addEventListener('click', toggleFavourite);
document.getElementById('favLink').addEventListener('click', showFavourites);
document.getElementById('favOverlay').addEventListener('click', (e) => {
    if (e.target === document.getElementById('favOverlay')) {
        hideFavourites();
    }
});
document.getElementById('favClose').addEventListener('click', hideFavourites);

(async () => {
    await loadCardsData();
    updateFavouriteCount();
    loadRandomCard();
})();
