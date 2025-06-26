import LandingHeader from "../component/LandingHeader";
import {useNavigate} from "react-router-dom";

function Landing() {
    const navigate = useNavigate();

    const handleSignInClick = () => {
         navigate("/register");
    }

    return (
        <>
        <LandingHeader />
        <div className="landing-container">
            <div><h1>Encuentra tu media <br/>musical</h1></div>
            <div>
                <button className="button" onClick={handleSignInClick}>Crear una cuenta</button>
            </div>
        </div>
        </>
    );
}

export default Landing;

