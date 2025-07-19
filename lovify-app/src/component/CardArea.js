import {useEffect, useState} from "react";
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import ClearIcon from '@mui/icons-material/Clear';
import FavoriteIcon from '@mui/icons-material/Favorite';
import {useConfig} from "../context/ConfigContext";

function CardArea() {
    const { apiUrl } = useConfig()
    const userID = sessionStorage.getItem("userID")
    const [user, setUser] = useState([]);
    const [slideIndex, setSlideIndex] = useState(0);
    const [notFound, setNotFound] = useState(true);

    const fetchRecommendations = async () => {
        try {
            const response = await fetch(`${apiUrl}/users/${userID}/recommendations`, {
                method: 'GET',
                credentials: 'include',
            })
            if (response.status === 404) {
                console.log(notFound)
                setNotFound(true);
            } else if (response.ok) {
                setNotFound(false);
                const data = await response.json();
                setUser(data);  // data = { recommended_user_id, similarity_score }
            } else {
                throw new Error("Error desconocido");
            }
        } catch (err) {
            console.error("Error al obtener recomendación:", err);
            setNotFound(true);
        }
    }

    useEffect(() => {
        fetchRecommendations();
    }, []);

    const currentSlide = user?.photos?.[slideIndex];

    if (!user || !currentSlide) return <p>Cargando...</p>;

    const nextSlide = () => {
        setSlideIndex((prev) =>
            prev + 1 < user.photos.length ? (prev + 1 % user.photos.length) : 0
        );
    };

    const prevSlide = () => {
        setSlideIndex((prev) =>
            prev - 1 >= 0 ? ((prev - 1 + user.photos.length) % user.photos.length) : user.photos.length - 1
        );
    };

    const like = (likeType) => {
        fetch(`${apiUrl}/users/${userID}/likes/${user.id}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    type: likeType,
                }),
            })
            .then((res) => res.json())
            .catch(console.error);

        fetchRecommendations();
        setSlideIndex(0); // reset al primer slide del siguiente user
    };

    if (notFound) return (
            <div className="error-card">
                <h1>No hay usuarios compatibles <br /> disponibles en este momento</h1>
            </div>
    );


    return (
        <div className="card-area">
            <div className="card">
               <img src={currentSlide} alt="current slice" className="card-image" />
                <div className="card__progress">
                    {user.photos.map((_, i) => (
                        <span
                            key={i}
                            className={`progress-dot ${i === slideIndex ? 'active' : ''}`}
                        />
                    ))}
                </div>
                <button className="arrow arrow--left" onClick={prevSlide}>
                    <ArrowBackIosNewIcon fontSize="small" style={{color: 'white'}}/>
                </button>
                <button className="arrow arrow--right" onClick={nextSlide}>
                    <ArrowForwardIosIcon fontSize="small" style={{color: 'white'}}/>
                </button>

                <div className="card__info">
                    <div className="card__info--title">
                        <h2>{user.name}</h2>
                        <p>{user.age} años</p>
                    </div>
                    {slideIndex === 0 && (
                        <>
                            <div className="top__track__header">
                                <p>🎵</p>
                                <p>Canción favorita</p>
                            </div>
                            <div className="top__track">
                                {/* eslint-disable-next-line jsx-a11y/img-redundant-alt */}
                                <img className="top_track_photo" src={user.top_tracks[0].album.cover}
                                     alt={"Top track Image"} height="80spx"/>
                                <div>
                                    <p className="top__track__title">{user.top_tracks[0].name}</p>
                                    <p>{user.top_tracks[0].artists[0]}</p>
                                </div>
                            </div>
                        </>
                    )
                    }
                    {slideIndex === 1 && (
                        <>
                            <div className="top__track__header">
                                <p>🎤</p>
                                <p>Artista Favorito</p>
                            </div>
                            <div className="top__track">
                                {/* eslint-disable-next-line jsx-a11y/img-redundant-alt */}
                                <img className="top_track_photo" src={user.top_artists[0].image} alt={"Top track Image"}
                                     height="80spx"/>
                                <div className="top__track_info">
                                    <p className="top__track__title">{user.top_artists[0].name}</p>
                                </div>
                            </div>
                        </>
                    )
                    }
                    {slideIndex === 2 && (
                        <>
                            <p className="bio">{user.bio}</p>
                        </>
                    )
                    }
                    {slideIndex === 3 && (
                        <>
                            <p className="favorite__artists">Artistas favoritos</p>
                            <div className="artist-row">
                                {user.top_artists.map((artist, index) => (
                                    <div key={index} className="artist">
                                        <img src={artist.image} alt={artist.name}/>
                                    </div>
                                ))}
                            </div>
                        </>
                    )
                    }
                </div>

                {currentSlide.extra && (
                    <img src={currentSlide.extra} className="artist-avatar" alt="extra"/>
                )}

                <div className="card__actions">
                    <button className="btn btn--no" onClick={() => like('dislike')}>
                        <ClearIcon fontSize="large" style={{color: '#FF3E3E'}}/>
                    </button>

                    <button className="btn btn--yes" onClick={() => like('like')}>
                        <FavoriteIcon fontSize="large" style={{color: '#00C851'}}/>
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CardArea;