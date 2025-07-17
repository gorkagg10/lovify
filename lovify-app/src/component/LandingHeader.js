import {useNavigate} from "react-router-dom";

function LandingHeader() {
    const navigate = useNavigate();

    const handleLoginClick = () => {
        navigate('/login')
    }

    return (
        <header className="header">
            <div className="header-container">
                <div className="login-wrapper">
                    <img src="/img.png" alt="Lovify Logo" className="lovify_logo"/>
                    <span className="brand">Lovify</span>
                </div>
            </div>
            <button className="login-button-header" onClick={handleLoginClick}>Inicia Sesión</button>
        </header>
    );
}

export default LandingHeader;