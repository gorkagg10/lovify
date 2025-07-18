import {Avatar} from "@mui/material";
import {useEffect, useState} from "react";
import {useConfig} from "../context/ConfigContext";
import {useNavigate} from "react-router-dom";

function Sidebar() {
    const [sidebarView, setSidebarView] = useState("matches");
    const [matches, setMatches] = useState([]);
    const { apiUrl } = useConfig()
    const userID = sessionStorage.getItem("userID")
    const navigate = useNavigate()

    const handleMessageClick = (matchID) => {
        navigate(`/messages/${matchID}`)
    }


    useEffect(() => {
        const fetchMatches = async () => {
            try {
                const response = await fetch(`${apiUrl}/users/${userID}/matches`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) throw new Error('Error al obtener los matches');

                const data = await response.json();
                setMatches(data);
            } catch (error) {
                console.error(error);
            }
        };

        fetchMatches();
    }, []);

    const conversations = [
        {
            id: "1",
            first_image: "/pexels-almadastudio-609549.jpg",
            name: "Fatima",
            age: 22,
            last_message: "Hola!"
        },
        {
            id: "2",
            first_image: "/pexels-kqpho-1921168.jpg",
            name: "Marina",
            age: 22,
            last_message: "Que guay!"
        }
    ]

    return (
        <aside className="sidebar">
            <header className="sidebar__header">
                <Avatar
                    alt="Gorka"
                    src="/gorka.png"
                    sx={{ width: 56, height: 56 }}
                />
                <span className="username">Gorka</span>
            </header>

            <nav className="sidebar__nav">
                <button className={`sidebar-btn ${sidebarView === "matches" ? "active" : ""}`}
                        onClick={() => setSidebarView("matches")}
                >
                    Matches
                </button>
                <button className={`sidebar-btn ${sidebarView === "messages" ? "active" : ""}`}
                        onClick={() => setSidebarView("messages")}
                >
                    Mensajes
                </button>
            </nav>

            {sidebarView === "matches" ? (
                <ul className="matches">
                    {matches.map((match) => (
                        <li key={match.id} className="match" onClick={() => handleMessageClick(match.match_id)}>
                            <Avatar
                                variant= "rounded"
                                className="match__avatar"
                                alt="Gorka"
                                src={match.first_image}
                            />
                            <span className="match__name">{match.name}</span>
                        </li>
                    ))}
                </ul>
                ): (
                <ul className="matches">
                    {conversations.map((convo) => (
                        <li key={convo.id} className="match">
                            <Avatar
                                variant="rounded"
                                className="match__avatar"
                                alt="Gorka"
                                src={convo.first_image}
                            />
                            <div className="conversation-info">
                                <span className="match__name">{convo.name}</span>
                                <p className="last-message">{convo.last_message}</p>
                            </div>
                        </li>
                    ))}
                </ul>
            )
            }
        </aside>
    );
}

export default Sidebar;