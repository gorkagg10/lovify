import {useState} from "react";
import {useNavigate} from "react-router-dom";
import {useConfig} from "../context/ConfigContext";

function CreateUser() {
    const { apiUrl } = useConfig();
    const [name, setName] = useState('');
    const [birthday, setBirthday] = useState('');
    const [gender, setGender] = useState('');
    const [sexualOrientation, setSexualOrientation] = useState('');
    const [description, setDescription] = useState('');
    const [message, setMessage] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const email = 'gorka.gonzalez.gomez@gmail.com';
            // localStorage.getItem('email')
        const formData = new FormData();

        if (email) {
            formData.append("email", email);
        }
        formData.append("name", name)
        formData.append("birthday", birthday);
        formData.append("gender", gender)
        formData.append("sexualOrientation", sexualOrientation)
        formData.append("description", description);

        const response = await fetch(`${apiUrl}/users`, {
            method:'POST',
            contentType:'multipart/form-data',
            body: formData,
            credentials: 'include'
        })

        if (response.status === 201) {
            const data = await response.json();
            console.log(data.id);
            sessionStorage.setItem('userID', data.id);
            navigate('/spotify-connect')
        } else {
            setMessage('Error creating user');
        }
    }

    return (
        <div className="create-user-container">
            <img src="logo.png" alt="Logo" className="logo"/>
            <form className="create-user-form" onSubmit={handleSubmit}>
                <div className="row">
                    <div className="input-wrapper">
                        <label>Nombre</label>
                        <input
                            id="name"
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            required
                        />
                    </div>
                    <div className="input-wrapper">
                        <label>Fecha de Nacimiento</label>
                        <input
                            id="birthday"
                            type="date"
                            value={birthday}
                            onChange={(e) => setBirthday(e.target.value)}
                            required
                        />
                    </div>
                </div>
                <div className="row">
                    <div className="input-wrapper">
                        <label htmlFor="genre">Género</label>
                        <select
                            id="genre"
                            name="genre"
                            value={gender}
                            onChange={(e) => setGender(e.target.value)}
                            required>
                            <option value="">Selecciona una opción</option>
                            <option value="male">Hombre</option>
                            <option value="female">Mujer</option>
                        </select>
                    </div>
                    <div className="input-wrapper">
                        <label htmlFor="sexualOrientation">Orientación sexual</label>
                        <select
                            id="sexualOrientation"
                            name="sexualOrientation"
                            value={sexualOrientation}
                            onChange={(e) => setSexualOrientation(e.target.value)}
                            required>
                            <option value="">Selecciona una opción</option>
                            <option value="heterosexual">Heterosexual</option>
                            <option value="homosexual">Homosexual</option>
                        </select>
                    </div>
                </div>
                <div className="row">
                    <div className="input-wrapper-description">
                        <label>Descripción</label>
                        <textarea
                            id="description"
                            name="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            rows="6"
                            required
                        />
                    </div>
                </div>
                <button type="submit" className="create-user-button">Crear Perfil</button>
            </form>
        </div>
    );
}

export default CreateUser;