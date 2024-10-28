function createCard(data,root) {
    let name = data["name"];
    let shortName = data["shortName"];
    let comments = data["amountComments"];
    let rating = data["rating"];
    let id = data["id"];

    let tempDiv = document.createElement('div');
    let activeStar = function (number) {
        return data["rating"] < number ? "star-gray" : "";
    }
    tempDiv.innerHTML = `{{template "teacher_card" .}}`;
    root.appendChild(tempDiv.firstElementChild);
    tempDiv.remove();
}

document.addEventListener('DOMContentLoaded', () => {
    const teacher_cards = document.getElementById('teacher_cards');
    let data = {"name" : "Ivana Tolic-Sapine","shortName":"TSI","amountComments":10,"rating":4,"id":1}
    let data2 = {"name" : "Ivana Tolic-Sapine","shortName":"TSI","amountComments":10,"rating":1,"id":2}
    createCard(data, teacher_cards);
    createCard(data2, teacher_cards);
})
document.getElementById("filterToggleButton").addEventListener("click", function() {
    const filterMenu = document.getElementById("filterMenu");
    filterMenu.classList.toggle("active");
});
