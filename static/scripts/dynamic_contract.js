let dynamic_contract = document.getElementById("contract-price");
let lang = document.getElementById("contract-lang");
let month = document.getElementById("contract-month");

let english = document.getElementById("contract-english");
let germany = document.getElementById("contract-germany");
let french = document.getElementById("contract-french");
let spanish = document.getElementById("contract-spanish");
let china = document.getElementById("contract-china");
let japan = document.getElementById("contract-japan");
let hindi = document.getElementById("contract-hindi");
let hebrew = document.getElementById("contract-hebrew");
let kazakh = document.getElementById("contract-kazakh");
let chuvash = document.getElementById("contract-chuvash");
let turkish = document.getElementById("contract-turkish");
let arabic = document.getElementById("contract-arabic");

const price_english = 1700;
const price_germany = 1500;
const price_french = 1400;
const price_spanish = 1600;
const price_china = 2000;
const price_japan = 1550;
const price_hindi = 1650;
const price_hebrew = 1750;
const price_kazakh = 1800;
const price_chuvash = 2100;
const price_turkish = 1680;
const price_arabic = 1900;




function getInfo(){
    let val = lang.value;
    if(val === "Английский язык") {
        dynamic_contract.setAttribute("value", month.value * price_english + "₽");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        english.classList.add("displayLanguage");
    }
    else if(val === "Немецкий язык") {
        dynamic_contract.setAttribute("value", month.value * price_germany + "₽");
        english.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        germany.classList.add("displayLanguage");
    }
    else if(val === "Французский язык"){
        dynamic_contract.setAttribute("value",month.value * price_french+"₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        french.classList.add("displayLanguage");
    }
    else if(val === "Испанский язык") {
        dynamic_contract.setAttribute("value", month.value * price_spanish + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        spanish.classList.add("displayLanguage");
    }
    else if(val === "Китайский язык") {
        dynamic_contract.setAttribute("value", month.value * price_china + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        china.classList.add("displayLanguage");
    }
    else if(val === "Японский язык") {
        dynamic_contract.setAttribute("value", month.value * price_japan + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        japan.classList.add("displayLanguage");
    }
    else if(val === "Хинди") {
        dynamic_contract.setAttribute("value", month.value * price_hindi + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        hindi.classList.add("displayLanguage");
    }
    else if(val === "Иврит") {
        dynamic_contract.setAttribute("value", month.value * price_hebrew + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        hebrew.classList.add("displayLanguage");
    }
    else if(val === "Казахский язык") {
        dynamic_contract.setAttribute("value", month.value * price_kazakh + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        kazakh.classList.add("displayLanguage");
    }
    else if(val === "Чувашский язык") {
        dynamic_contract.setAttribute("value", month.value * price_chuvash + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        chuvash.classList.add("displayLanguage");
    }
    else if(val === "Турецкий язык") {
        dynamic_contract.setAttribute("value", month.value * price_turkish + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        arabic.classList.remove("displayLanguage");
        turkish.classList.add("displayLanguage");
    }
    else if(val === "Арабский язык") {
        dynamic_contract.setAttribute("value", month.value * price_arabic + "₽");
        english.classList.remove("displayLanguage");
        germany.classList.remove("displayLanguage");
        french.classList.remove("displayLanguage");
        spanish.classList.remove("displayLanguage");
        china.classList.remove("displayLanguage");
        japan.classList.remove("displayLanguage");
        hebrew.classList.remove("displayLanguage");
        hindi.classList.remove("displayLanguage");
        kazakh.classList.remove("displayLanguage");
        chuvash.classList.remove("displayLanguage");
        turkish.classList.remove("displayLanguage");
        arabic.classList.add("displayLanguage");
    }
}

window.onload = function () {

}

