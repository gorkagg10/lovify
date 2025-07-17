import {useNavigate} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function SpotifyConnect() {
    const { apiUrl } = useConfig()
    console.log(apiUrl)

    const handleLogin= () => {
        const userID = sessionStorage.getItem("userID");
        window.location.href = `${apiUrl}/users/${userID}/login/spotify`;
    }
    return (
        <div className="spotify-connect-container">
            <img src="logo.png" alt="Logo" className="logo"/>
            <button type="submit" className="spotify-connect-button" onClick={handleLogin}>Conectar con perfil de Spotify</button>
        </div>
    )
}

export default SpotifyConnect;