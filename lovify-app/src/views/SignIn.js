import {useState} from "react";
import {useNavigate} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function SignIn() {
    const { apiUrl } = useConfig();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("email", email);
        formData.append("password", password);

        const response = await fetch(`${apiUrl}/auth/register`, {
            method:'POST',
            contentType:'multipart/form-data',
            body: formData,
        });

        if (response.status === 201) {
            navigate('/login')
        } else {
            setMessage('Error al iniciar sesión ❌')
        }
    }

    return (
        <div className="sigin-container">
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

                <button className="signin-button" type="submit">Crear cuenta</button>

                {message && <p className="message">{message}</p>}
            </form>
        </div>

    );
}

export default SignIn;