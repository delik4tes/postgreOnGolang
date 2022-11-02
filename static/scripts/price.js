let price = document.getElementById("contract-price");
let lang = document.getElementById("contract-lang");
let month = document.getElementById("contract-month");

const price_english = 1700;
const price_germany = 1500;
const price_franch = 1400;
const price_spanish = 1600;
const price_china = 2000;
const price_japan = 1550;
const price_hindi = 1650;
const price_ivrit = 1750;
const price_kazah = 1800;
const price_chuvas = 2100;
const price_turkish = 1680;
const price_arabic = 1900;

function getprice(){
    let val = lang.value;
    if(val === "Английский язык")
        price.setAttribute("value",month.value * price_english+"₽");
    else if(val === "Немецкий язык")
        price.setAttribute("value",month.value * price_germany+"₽");
    else if(val === "Французский язык")
        price.setAttribute("value",month.value * price_franch+"₽");
    else if(val === "Испанский язык")
        price.setAttribute("value",month.value * price_spanish+"₽");
    else if(val === "Китайский язык")
        price.setAttribute("value",month.value * price_china+"₽");
    else if(val === "Японский язык")
        price.setAttribute("value",month.value * price_japan+"₽");
    else if(val === "Хинди")
        price.setAttribute("value",month.value * price_hindi+"₽");
    else if(val === "Иврит")
        price.setAttribute("value",month.value * price_ivrit+"₽");
    else if(val === "Казахский язык")
        price.setAttribute("value",month.value*price_kazah+"₽");
    else if(val === "Чувашский язык")
        price.setAttribute("value",month.value*price_chuvas+"₽");
    else if(val === "Турецкий язык")
        price.setAttribute("value",month.value*price_turkish+"₽");
    else if(val === "Арабский язык")
        price.setAttribute("value",month.value*price_arabic+"₽");
}

