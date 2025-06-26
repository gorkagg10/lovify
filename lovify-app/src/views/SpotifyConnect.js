import {useNavigate} from "react-router-dom";

function SpotifyConnect() {
    const navigate = useNavigate();

    const handleLogin= () => {
        const userID = sessionStorage.getItem("userID");
        window.location.href = `http://localhost:8080/users/${userID}/login/spotify`;
    }
    return (
        <div className="spotify-connect-container">
            <img src="logo.png" alt="Logo" className="logo"/>
            <button type="submit" onClick={handleLogin}>Conectar con perfil de Spotify</button>
        </div>
    )
}

export default SpotifyConnect;