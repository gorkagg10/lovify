import {Avatar} from "@mui/material";
import {useState} from "react";

function Sidebar() {
    const [sidebarView, setSidebarView] = useState("matches");
    const matches = [
        {
            id: "1",
            image: "/pexels-hannah-nelson-390257-1065084.jpg",
            age: 22,
            name: "Claudia",
        },
        {
            id: "2",
            image: "/pexels-pixabay-415829.jpg",
            age: 18,
            name: "Julia",
        }
    ];
    const conversations = [
        {
            id: "1",
            image: "/pexels-almadastudio-609549.jpg",
            name: "Fatima",
            age: 22,
            last_message: "Hola!"
        },
        {
            id: "2",
            image: "/pexels-kqpho-1921168.jpg",
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
                        <li key={match.id} className="match">
                            <Avatar
                                variant= "rounded"
                                className="match__avatar"
                                alt="Gorka"
                                src={match.image}
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
                                src={convo.image}
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