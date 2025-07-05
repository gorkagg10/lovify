import CardArea from "../component/CardArea";
import Sidebar from "../component/Sidebar";

function MainPage() {
    return (
        <div className="main_page">
            <Sidebar/>
            <CardArea/>
        </div>
    )
}

export default MainPage;