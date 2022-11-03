let experience = document.getElementById("experience");
let teacher_lang = document.getElementById("teacher_lang");
let address = document.getElementById("address");
let phone = document.getElementById("phone");
let insertAddress = document.getElementById("insert-address");

let teacher = document.getElementById("teacher");
let student = document.getElementById("student");
let admin = document.getElementById("admin");

function change (){
    if (teacher.checked){
        experience.classList.add("displayBlock");
        teacher_lang.classList.add("displayBlock");
        address.classList.add("displayBlock");
        phone.classList.remove("displayBlock");
        insertAddress.classList.remove("displayBlock");
    }
    if (student.checked){
        experience.classList.remove("displayBlock");
        teacher_lang.classList.remove("displayBlock");
        address.classList.add("displayBlock");
        phone.classList.add("displayBlock");
        insertAddress.classList.remove("displayBlock")
    }
    if(admin.checked){
        experience.classList.remove("displayBlock");
        teacher_lang.classList.remove("displayBlock");
        address.classList.remove("displayBlock");
        phone.classList.remove("displayBlock");
        insertAddress.classList.add("displayBlock")
    }
}

window.onload = function () {
    student.checked = true;
}