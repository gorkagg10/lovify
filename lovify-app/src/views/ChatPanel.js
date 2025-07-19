import {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function ChatPanel() {
    const [messages, setMessages] = useState([]);
    const { apiUrl } = useConfig()
    const userID = sessionStorage.getItem("userID")
    const [receiver, setReceiver] = useState([]);
    const { matchId } = useParams();

    const fetchReceiverInfo = async () => {
        const receiverUserID = sessionStorage.getItem("receiverID")
        try {
            const response = await fetch(`${apiUrl}/users/${receiverUserID}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) throw new Error('Error al obtener los matches');

            const data = await response.json();
            setReceiver(data);
        } catch (error) {
            console.error(error);
        }
    }

    const fetchMessages = async () => {
        try {
            const response = await fetch(`${apiUrl}/users/${userID}/messages/${matchId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) throw new Error('Error al obtener los matches');

            const data = await response.json();
            setMessages(data);
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        fetchReceiverInfo();
        fetchMessages();
    }, []);

    const [content, setContent] = useState("");
    const [sending, setSending] = useState(false);
    const [error, setError] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSending(true);
        setError("");

        try {
            const res = await fetch(`${apiUrl}/users/${userID}/messages/${matchId}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json", // JWT
                },
                body: JSON.stringify({ content: content }),
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.message || "Error al enviar mensaje");
            }

            setContent("");

            fetchMessages();
        } catch (err) {
            setError(err.message);
        } finally {
            setSending(false);
        }
    };

    return (
        <div className="chat-panel">
            <div className="chat-panel__header">
                <img className="chat-panel__logo" src={receiver.photos?.[0]} alt={receiver.name}/>
                <span>{receiver.name}, {receiver.age}</span>
            </div>

            <div className="chat-messages">
                {messages.map((msg, i) => (
                    <div
                        key={i}
                        className={`chat-bubble ${msg.from_user_id === userID ? 'from-me' : 'from-them'}`}
                    >
                        {msg.content}
                    </div>
                ))}
            </div>
            <div className="chat-input">
                <input
                    type="text"
                    placeholder="Escribe un mensaje..."
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    disabled={sending}
                    required
                    style={{ flexGrow: 1 }}
                />
                <button type="submit" disabled={sending || !content.trim()} onClick={handleSubmit}>ENVIAR</button>
            </div>
        </div>
    )
}

export default ChatPanel;