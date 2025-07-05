import {useState} from "react";
import {useConfig} from "../context/ConfigContext";

const MAX_IMAGES = 6;

function ProfilePhotoUploader() {
    const { apiUrl } = useConfig()
    const [photos, setPhotos] = useState([]);

    const handleRemove = (index) => {
        const updated = [...photos];
        updated.splice(index, 1);
        setPhotos(updated);
    };

    const handleAdd = (e) => {
        const file = e.target.files[0];
        if (!file) return;
        const url = URL.createObjectURL(file);
        setPhotos((prev) => [...prev, url]);
    };

    const handleSave = async () => {
        const userID = sessionStorage.getItem("userID");

        const formData = new FormData();

        photos.forEach((photo, i) => {
            formData.append('photos[]', photo);
        })

        const response = await fetch(`${apiUrl}/users/${userID}/photos`, {
            method:'POST',
            body: formData,
        })

        console.log(response)
    }

    return (
        <div className="photo-uploader-container">
            <h2>FOTOS DE PERFIL</h2>
            <div className="photo-grid">
                {Array.from({length: MAX_IMAGES}).map((_, i) => (
                    <div className="photo-slot" key={i}>
                        {photos[i] ? (
                            <>
                                <img src={photos[i]} alt={`Foto ${i}`}/>
                                <button className="remove-btn" onClick={() => handleRemove(i)}>
                                    ×
                                </button>
                            </>
                        ) : (
                            <label className="add-btn">
                                +
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleAdd}
                                    hidden
                                />
                            </label>
                        )}
                    </div>
                ))}
            </div>
            <button className="save-btn" onClick={handleSave}>Guardar</button>
        </div>
    );
}

export default ProfilePhotoUploader;