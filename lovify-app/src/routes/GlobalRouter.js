import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Landing from "../views/Landing";
import Login from "../views/Login";
import SignIn from "../views/SignIn";
import CreateUser from "../views/CreateUser";
import SpotifyConnect from "../views/SpotifyConnect";
import MainPage from "../views/MainPage";
import Messages from "../views/Messages";
import ProfilePhotos from "../views/ProfilePhotos";


function GlobalRouter() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Landing />}/>
                <Route path="/login" element={<Login />}/>
                <Route path="/register" element={<SignIn />}/>
                <Route path="/users" element={<CreateUser />}/>
                <Route path="/app/profile/photos" element={<ProfilePhotos />}/>
                <Route path="/spotify-connect" element={<SpotifyConnect />}/>
                <Route path="/app" element={<MainPage />}/>
                <Route path="/messages/:matchId" element={<Messages />}/>
            </Routes>
        </BrowserRouter>
    );
}

export default GlobalRouter;