import {useState} from "react";
import {useNavigate} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function Login() {
    const { apiUrl } = useConfig()
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("email", email);
        formData.append("password", password);

        const response = await fetch(`${apiUrl}/auth/login`, {
            method:'POST',
            contentType:'multipart/form-data',
            body: formData,
            credentials: 'include'
        });

        if (response.ok) {
            sessionStorage.setItem('email', email)
            const data = await response.json();
            console.log(data);
            if (data.is_profile_connected) {
                sessionStorage.setItem('userID', data.profile_id);
                navigate('/app')
            }
            else {
                navigate('/users')
            }
        } else {
            setMessage('Error logging in!');
        }
    }

    return (
        <div className="login-container">
            <img src="logo.png" alt="Logo" className="logo"/>
            <form className="form" onSubmit={handleSubmit}>
                <div className="input-wrapper">
                    <label>Dirección de Correo</label>
                    <input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>
                <div className="input-wrapper">
                    <label>Contraseña</label>
                    <input
                        id="password"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>

                <button className="login-button" type="submit">Iniciar Sesión</button>

                {message && <p className="login-message">{message}</p>}
            </form>
        </div>

    );
}

export default Login;