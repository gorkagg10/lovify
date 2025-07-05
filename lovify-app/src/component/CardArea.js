import {useEffect, useState} from "react";
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import ClearIcon from '@mui/icons-material/Clear';
import FavoriteIcon from '@mui/icons-material/Favorite';

function CardArea() {
    const [users, setUsers] = useState([]);
    const [index, setIndex] = useState(0);
    const [slideIndex, setSlideIndex] = useState(0);

    useEffect(() => {
        // Simulación de llamada a la API
        const data = [{
            id: "1",
            slides: [
                {
                    image: "gorka.png",
                    title: "Sarah",
                    subtitle: "22 años",
                    extra: null
                },
                {
                    image: "gorka.png",
                    title: "Sarah",
                    subtitle: "22 años",
                    extra: null
                },
            ],
        }]
        setUsers(data);
    }, []);

    const user = users[index];
    const currentSlide = user?.slides?.[slideIndex];

    if (!user || !currentSlide) return <p>Cargando...</p>;

    const nextSlide = () => {
        setSlideIndex((prev) =>
            prev + 1 < user.slides.length ? prev + 1 : 0
        );
    };

    const prevSlide = () => {
        setSlideIndex((prev) =>
            prev - 1 >= 0 ? prev - 1 : user.slides.length - 1
        );
    };

    const nextUser = () => {
        setIndex((prev) => (prev + 1) % users.length);
        setSlideIndex(0); // reset al primer slide del siguiente user
    };

    return (
        <div className="card-area">
            <div className="card" style={{ backgroundImage: `url(${currentSlide.image})` }}>
                <div className="card__progress">
                    {user.slides.map((_, i) => (
                        <span
                            key={i}
                            className={`progress-dot ${i === slideIndex ? 'active' : ''}`}
                        />
                    ))}
                </div>
                <button className="arrow arrow--left" onClick={prevSlide}>
                    <ArrowBackIosNewIcon fontSize="small" style={{ color: 'white' }} />
                </button>
                <button className="arrow arrow--right" onClick={nextSlide}>
                    <ArrowForwardIosIcon fontSize="small" style={{ color: 'white' }} />
                </button>

                <div className="card__info">
                    <h2>{currentSlide.title}</h2>
                    <p>{currentSlide.subtitle}</p>

                    {currentSlide.extra && (
                        <img src={currentSlide.extra} className="artist-avatar" alt="extra" />
                    )}
                </div>

                <div className="card__actions">
                    <button className="btn btn--no">
                        <ClearIcon fontSize="large" style={{ color: '#FF3E3E' }} />
                    </button>

                    <button className="btn btn--yes">
                        <FavoriteIcon fontSize="large" style={{ color: '#00C851' }} />
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CardArea;