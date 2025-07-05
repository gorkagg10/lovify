function ChatPanel() {
    const selectedUser = {
        id: "2",
        image: "gorka.png",
        name: "Gorka",
        age: 22,
        last_message: "hola k ase 2 me llamo gorkan jjjjjjjjjjjjj"
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
            text: "¡Hola Júlia! Muy bien, ¿y tú?",
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
            text: "Sin duda *Starboy*. ¿Y la tuya?",
            timestamp: "2025-07-03T19:33:00Z",
        },]

    return (
        <div className="chat-panel">
            {/*<div className="chat-panel__header">
                <img className="chat-panel__logo" src={selectedUser.image} alt={selectedUser.name}/>
                <span>{selectedUser.name}, {selectedUser.age}</span>
            </div>*/}

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
        </div>
    )
}

export default ChatPanel;