import {Avatar} from "@mui/material";
import {useState} from "react";

function Sidebar() {
    const [sidebarView, setSidebarView] = useState("matches");
    const matches = [
        {
            id: "1",
            image: "gorka.png",
            age: 22,
            name: "Gorka",
        },
        {
            id: "2",
            image: "gorka.png",
            age: 18,
            name: "Gorka",
        }
    ];
    const conversations = [
        {
            id: "1",
            image: "gorka.png",
            name: "Rahima",
            age: 22,
            last_message: "hola k ase"
        },
        {
            id: "2",
            image: "gorka.png",
            name: "Gorka",
            age: 22,
            last_message: "hola k ase 2 me llamo gorkan jjjjjjjjjjjjj"
        }
    ]

    return (
        <aside className="sidebar">
            <header className="sidebar__header">
                <Avatar
                    alt="Gorka"
                    src="gorka.png"
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