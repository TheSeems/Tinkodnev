const ranks = {
    0: "<span class='blacklist'>Черный список</span>",
    1: "<span class='student'>Ученик</span>",
    2: "<span class='teacher'>Преподаватель</span>",
    3: "<span class='organizer'>Организатор</span>"
};

function httpGetAsync(url, method, callback) {
    let xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState === 4 && xmlHttp.status === 200)
            callback(xmlHttp.responseText);
        else if (xmlHttp.status !== 200 && xmlHttp.status !== 0) {
            alert("Не могу ничего достать!");
            console.log(xmlHttp.readyState, xmlHttp.status)
        }
    };
    xmlHttp.open(method, url, true); // true for asynchronous
    xmlHttp.send(null);
}

function load(id) {
    httpGetAsync("/api/get?id=" + id, "GET", function (json) {
        let parsed = JSON.parse(json);
        console.log(parsed);

        if (parsed["success"] === false) {
            if (parsed["error"] === "Not found")
                document.getElementById("main").innerText = "Никого не найдено!";
            else
                document.getElementById("main").innerText = "ОШИБКА: " + parsed["error"];
            return
        }

        document.getElementById("main").innerText = "Пользователь под айди " + id + ":";
        document.getElementById("info").style.display = "block";
        document.getElementById("user_firstName").innerHTML = parsed["member"]["first_name"];
        document.getElementById("user_secondName").innerHTML = parsed["member"]["second_name"];
        document.getElementById("user_status").innerHTML = ranks[parsed["member"]["status"]];
    });
}