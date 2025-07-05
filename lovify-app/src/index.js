import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import './styles/landing.css'
import './styles/landing_header.css';
import './styles/login.css';
import './styles/signin.css';
import './styles/create_user.css';
import './styles/spotify_connect.css'
import './styles/card_area.css'
import './styles/sidebar.css'
import './styles/main_page.css';
import './styles/messages.css';
import './styles/chat_panel.css';
import './styles/profile_photo_uploader.css';
import './styles/profile_photos.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import {ConfigProvider} from "./context/ConfigContext";

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
      <ConfigProvider>
          <App />
      </ConfigProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
