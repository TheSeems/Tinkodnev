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
        if (parsed["member"]["patronymic"] !== " ")
            document.getElementsByClassName("user_patronymic")[0].innerHTML = parsed["member"]["patronymic"];
        else
            document.getElementsByClassName("patronymic")[0].remove();

        document.getElementsByClassName("user_firstName")[0].innerHTML = parsed["member"]["first_name"];
        document.getElementsByClassName("user_secondName")[0].innerHTML = parsed["member"]["second_name"];
        document.getElementsByClassName("user_status")[0].innerHTML = ranks[parsed["member"]["status"]];
        document.getElementById("info").style.display = "block";
    });
}

function find(query) {
    httpGetAsync("/api/search?query=" + query, "GET", function (json) {
        let parsed = JSON.parse(json);
        console.log(parsed);

        let ul = document.getElementById("list");
        ul.innerHTML = "";
        if (parsed["success"] === false) {
            if (parsed["error"] === "Not found")
                document.getElementById("main").innerText = "Никого не найдено!";
            else
                document.getElementById("main").innerText = "ОШИБКА: " + parsed["error"];
        } else {
            document.getElementById("main").innerText = "Найдено пользователей: " + parsed["members"].length;
            for (let loh in parsed["members"]) {
                let i = parsed["members"][loh];
                let li = document.createElement("li");
                let element = document.createElement("user");
                element.innerHTML = "<span class = 'user_firstName'>" + i["first_name"] + "</span> <span class = 'user_secondName'>"
                    + i["second_name"] + "</span> <span class='user_status'>" + ranks[i["status"]] + "</span> (<a href='/view?id=" + i["id"] + "'>" + i["id"] + "</a>)";
                li.appendChild(element);
                ul.appendChild(li);
            }
        }
    });
}