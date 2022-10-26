let experience = document.getElementById("experience");
let teacher_lang = document.getElementById("teacher_lang");
let salary = document.getElementById("salary");
let address = document.getElementById("address");
let phone = document.getElementById("phone");

let teacher = document.getElementById("teacher");
let student = document.getElementById("student");
let admin = document.getElementById("admin");

function change (){
    if (teacher.checked){
        experience.classList.add("displayBlock");
        teacher_lang.classList.add("displayBlock");
        // salary.classList.add("displayBlock");
        address.classList.add("displayBlock");
        phone.classList.remove("displayBlock");
    }
    if (student.checked){
        experience.classList.remove("displayBlock");
        teacher_lang.classList.remove("displayBlock");
        // salary.classList.remove("displayBlock");
        address.classList.add("displayBlock");
        phone.classList.add("displayBlock");
    }
    if(admin.checked){
        experience.classList.remove("displayBlock");
        teacher_lang.classList.remove("displayBlock");
        // salary.classList.add("displayBlock");
        address.classList.add("displayBlock");
        phone.classList.remove("displayBlock");
    }
}

window.onload = function () {
    student.checked = true;
}