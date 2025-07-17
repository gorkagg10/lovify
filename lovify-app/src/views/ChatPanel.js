import {useState} from "react";

function ChatPanel() {
    const selectedUser = {
        id: "2",
        image: "/pexels-kqpho-1921168.jpg",
        name: "Marina",
        age: 22,
        timestamp: "2025-07-02T19:30:00Z",
    }

    const messages = [
        {
            id: 1,
            sender: "them", // mensaje recibido
            text: "¡Hola! ¿Cómo estás?",
            timestamp: "2025-07-03T19:30:00Z",
        },
        {
            id: 2,
            sender: "me", // mensaje enviado por ti
            text: "¡Hola Marina! Muy bien, ¿y tú?",
            timestamp: "2025-07-03T19:31:00Z",
        },
        {
            id: 3,
            sender: "them",
            text: "Genial. Me encantó tu perfil musical 😊",
            timestamp: "2025-07-03T19:32:00Z",
        },
        {
            id: 4,
            sender: "me",
            text: "Gracias 😄 ¿Cuál es tu canción favorita de The Weeknd?",
            timestamp: "2025-07-03T19:32:30Z",
        },
        {
            id: 5,
            sender: "them",
            text: "Sin duda Starboy. ¿Y la tuya?",
            timestamp: "2025-07-03T19:33:00Z",
        },
        {
            id: 6,
            sender: "me",
            text: "Que guay!",
            timestamp: "2025-07-03T19:34:00Z",
        },
    ]

    return (
        <div className="chat-panel">
            <div className="chat-panel__header">
                <img className="chat-panel__logo" src={selectedUser.image} alt={selectedUser.name}/>
                <span>{selectedUser.name}, {selectedUser.age}</span>
            </div>

            <div className="chat-messages">
                {messages.map((msg, i) => (
                    <div
                        key={i}
                        className={`chat-bubble ${msg.sender === 'me' ? 'from-me' : 'from-them'}`}
                    >
                        {msg.text}
                    </div>
                ))}
            </div>
            <div className="chat-input">
                <input
                    type="text"
                    placeholder="Escribe un mensaje..."
                />
                <button>ENVIAR</button>
            </div>
        </div>
    )
}

export default ChatPanel;