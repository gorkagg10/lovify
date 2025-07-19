import {Avatar} from "@mui/material";
import {useEffect, useState} from "react";
import {useConfig} from "../context/ConfigContext";
import {useNavigate} from "react-router-dom";

function Sidebar() {
    const [sidebarView, setSidebarView] = useState("matches");
    const [matches, setMatches] = useState([]);
    const [conversations, setConversations] = useState("");
    const { apiUrl } = useConfig()
    const userID = sessionStorage.getItem("userID")
    const navigate = useNavigate()
    const [user, setUser] = useState([]);

    const handleMessageClick = (matchID, receiverID) => {
        sessionStorage.setItem("receiverID", receiverID)
        navigate(`/messages/${matchID}`)
    }

    useEffect(() => {
        const fetchUserInfo = async () => {
            try {
                const response = await fetch(`${apiUrl}/users/${userID}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });

                if (!response.ok) throw new Error('Error al obtener los matches');

                const data = await response.json();
                setUser(data);
            } catch (error) {
                console.error(error);
            }
        }

        const fetchMatches = async () => {
            try {
                const response = await fetch(`${apiUrl}/users/${userID}/matches`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });

                if (!response.ok) throw new Error('Error al obtener los matches');

                const data = await response.json();
                setMatches(data);
            } catch (error) {
                console.error(error);
            }
        };

        const fetchConversations = async () => {
            try {
                const response = await fetch(`${apiUrl}/users/${userID}/conversations`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });

                if (!response.ok) throw new Error('Error al obtener las conversacions');

                const data = await response.json();
                setConversations(data);
            } catch (error) {
                console.error(error);
            }
        };

        fetchUserInfo().then();
        fetchMatches().then();
        fetchConversations().then();
    }, [apiUrl, userID]);

    return (
        <aside className="sidebar">
            <div className="sidebar-content">
                <header className="sidebar__header">
                    <Avatar
                        alt={user.name}
                        src={user.photos?.[0]}
                        sx={{width: 56, height: 56}}
                    />
                    <span className="username">{user.name}</span>
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
                            <li key={match.id} className="match"
                                onClick={() => handleMessageClick(match.match_id, match.user_id)}>
                                <Avatar
                                    variant="rounded"
                                    className="match__avatar"
                                    alt="Gorka"
                                    src={match.first_image}
                                />
                                <span className="match__name">{match.name}</span>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <ul className="matches">
                        {conversations.map((convo) => (
                            <li key={convo.id} className="match"
                                onClick={() => handleMessageClick(convo.match_id, convo.user_id)}>
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
            </div>
                <button
                    onClick={() => {
                        localStorage.removeItem("token");
                        window.location.href = "/login";
                    }}
                    className="logout-button"
                >
                    Cerrar sesión
                </button>
        </aside>
);
}

export default Sidebar;