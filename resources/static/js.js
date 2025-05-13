function receiveMessage(event) {
    const chatDiv = document.getElementById("chat");
    const msg = JSON.parse(event.data);
    const currentUser = document.getElementById("username").value || "Noname";

    const msgContainer = document.createElement("div");
    msgContainer.className = "d-flex flex-column " + (msg.username === currentUser ? "align-items-end" : "align-items-start");

    const bubble = document.createElement("div");
    bubble.className = "message " + (msg.username === currentUser ? "self" : "other");

    const nameEl = document.createElement("div");
    nameEl.className = "username";
    nameEl.textContent = msg.username;

    const textEl = document.createElement("div");
    textEl.textContent = msg.message;

    bubble.appendChild(nameEl);
    bubble.appendChild(textEl);
    msgContainer.appendChild(bubble);
    chatDiv.appendChild(msgContainer);

    chatDiv.scrollTop = chatDiv.scrollHeight;
};

function sendMessage() {
    const username = document.getElementById("username").value || "Noname";
    const message = document.getElementById("message").value;
    if (message.trim() === "") return;
    ws.send(JSON.stringify({ username, message }));
    document.getElementById("message").value = "";
}


function loadRecentMessages() {
    fetch("/api/lastMessages")
        .then(response => response.json())
        .then(data => {
            const messages = data.messages;
            const currentUser = document.getElementById("username").value || "Noname";
            const chatDiv = document.getElementById("chat");

            messages.reverse();

            for (const msg of messages) {
                const msgContainer = document.createElement("div");
                msgContainer.className = "d-flex flex-column " + (msg.username === currentUser ? "align-items-end" : "align-items-start");

                const bubble = document.createElement("div");
                bubble.className = "message " + (msg.username === currentUser ? "self" : "other");

                const nameEl = document.createElement("div");
                nameEl.className = "username";
                nameEl.textContent = msg.username;

                const textEl = document.createElement("div");
                textEl.textContent = msg.message;

                bubble.appendChild(nameEl);
                bubble.appendChild(textEl);
                msgContainer.appendChild(bubble);
                chatDiv.appendChild(msgContainer);
            }

            chatDiv.scrollTop = chatDiv.scrollHeight;
        })
        .catch(err => console.error("Ошибка загрузки сообщений:", err));
}

