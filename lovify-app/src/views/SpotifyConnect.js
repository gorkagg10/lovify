import {useNavigate} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function SpotifyConnect() {
    const { apiUrl } = useConfig()

    const handleLogin= () => {
        const userID = sessionStorage.getItem("userID");
        window.location.href = `${apiUrl}/users/${userID}/login/spotify`;
    }
    return (
        <div className="spotify-connect-container">
            <img src="logo.png" alt="Logo" className="logo"/>
            <button type="submit" onClick={handleLogin}>Conectar con perfil de Spotify</button>
        </div>
    )
}

export default SpotifyConnect;