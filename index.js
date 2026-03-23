const CARDS_JSON = "cards.json";
const BASE_URL = 'https://nlr.ru/e-case3/sc2.php/web_gak/gc';
const MAX_ID = 133781;

let CARDS = {};

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
    const id = getRandomId();
    const cardNumber = getRandomCardNumber(id);
    const img = document.getElementById('card');
    img.src = `${BASE_URL}/${id}/${cardNumber}`;
    img.alt = `Card ${id}`;
}

document.getElementById('newCard').addEventListener('click', loadCard);

(async () => {
    await loadCardsData();
    loadCard();
})();
