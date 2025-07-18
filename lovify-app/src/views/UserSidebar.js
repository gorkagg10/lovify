import {useState} from "react";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";

function UserSidebar() {
    const [slideIndex, setSlideIndex] = useState(0);

    const user = {
        id: "1",
        name: "Sarah",
        age: "22",
        bio: "Soy una apasionada de la escalada",
        photos: [
            "pexels-olly-774909.jpg",
            "pexels-olly-733872.jpg",
            "pexels-jill-wellington-1638660-40192.jpg",
            "img_1.png",
        ],
        top_tracks: [
            {
                name: 'ACHO PR',
                album: {
                    name: 'nadie sabe lo que va a pasar mañana',
                    type: 'album',
                    cover: 'https://i.scdn.co/image/ab67616d0000b2732ea1f035463d11e1fc3b193d',
                },
                artists: [ 'Bad Bunny', 'Arcángel', 'De La Ghetto', 'Ñengo Flow' ]
            },
            {
                name: "I Can't Wait To Get There",
                album: {
                    name: 'Hurry Up Tomorrow',
                    type: 'album',
                    cover: 'https://i.scdn.co/image/ab67616d0000b273982320da137d0de34410df61',
                },
                artists: [ 'The Weeknd' ]
            },
            {
                name: 'EVIL J0RDAN',
                album: {
                    name: 'MUSIC',
                    type: 'album',
                    cover: 'https://i.scdn.co/image/ab67616d0000b2736b219c8d8462bfe254a20469'
                },
                artists: [ 'Playboi Carti' ]
            },
            {
                name: 'Wake Me Up (feat. Justice)',
                album: {
                    name: 'Hurry Up Tomorrow',
                    type: 'album',
                    cover: 'https://i.scdn.co/image/ab67616d0000b2734aa8b835e86c68b78ee147d6',
                },
                artists: [ 'The Weeknd', 'Justice' ]
            },
            {
                name: 'Johnny Glamour',
                album: {
                    name: 'DAISY',
                    type: 'album',
                    cover: 'https://i.scdn.co/image/ab67616d0000b2735f74cb93d0cda9c953648265',
                },
                artists: [ 'rusowsky', 'Las Ketchup' ]
            },
        ],
        top_artist: [
            {
                name: 'The Weeknd',
                genres: [],
                image:  'https://i.scdn.co/image/ab6761610000e5eb9e528993a2820267b97f6aae',
            },
            {
                name: 'Future',
                genres: [ 'rap' ],
                image: 'https://i.scdn.co/image/ab6761610000e5eb7565b356bc9d9394eefa2ccb',
            },
            {
                name: 'DELLAFUENTE',
                genres: [ 'flamenco urbano', 'flamenco' ],
                image: 'https://i.scdn.co/image/ab6761610000e5eb8bbfd0d545ef15b909ffe49e',
            },
            {
                name: 'Bad Bunny',
                genres: [ 'reggaeton', 'trap latino', 'urbano latino', 'latin' ],
                image: 'https://i.scdn.co/image/ab6761610000e5eb81f47f44084e0a09b5f0fa13',
            },
            {
                name: 'Brent Faiyaz',
                genres: [ 'r&b' ],
                image: 'https://i.scdn.co/image/ab6761610000e5eb2af1d912483f27af21eba49f',
            },
        ],
    }
    const currentSlide = user?.photos?.[slideIndex];

    const nextSlide = () => {
        setSlideIndex((prev) =>
            prev + 1 < user.photos.length ? (prev + 1 % user.photos.length ): 0
        );
    };

    const prevSlide = () => {
        setSlideIndex((prev) =>
            prev - 1 >= 0 ? ((prev - 1 + user.photos.length) % user.photos.length) : user.photos.length - 1
        );
    };

    return (
        <div className="card-area">
            <div className="card" style={{backgroundImage: `url(${currentSlide})`}}>
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
            </div>
        </div>
    )
}

export default UserSidebar;