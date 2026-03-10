const BASE_URL = 'https://nlr.ru/e-case3/sc2.php/web_gak/gc';
const MAX_ID = 133781;

function getRandomId() {
    return Math.floor(Math.random() * MAX_ID) + 1;
}

function loadCard() {
    const id = getRandomId();
    const img = document.getElementById('card');
    img.src = `${BASE_URL}/${id}/1`;
    img.alt = `Card ${id}`;
}

document.getElementById('newCard').addEventListener('click', loadCard);

loadCard();
