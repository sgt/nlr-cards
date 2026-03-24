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

function loadCard() {
    currentId = getRandomId();
    currentCardNumber = getRandomCardNumber(currentId);
    const img = document.getElementById('card');
    img.src = `${BASE_URL}/${currentId}/${currentCardNumber}`;
    img.alt = `Card ${currentId}/${currentCardNumber}`;
    updateFavouriteButton();
}

function getFavouriteKey(id, cardNumber) {
    return `fav_${id}_${cardNumber}`;
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
}

function updateFavouriteButton() {
    const btn = document.getElementById('favourite');
    btn.textContent = isFavourited(currentId, currentCardNumber) ? '🤍' : '❤️';
}

document.getElementById('newCard').addEventListener('click', loadCard);
document.getElementById('favourite').addEventListener('click', toggleFavourite);

(async () => {
    await loadCardsData();
    loadCard();
})();
