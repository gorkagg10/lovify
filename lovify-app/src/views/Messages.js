import Sidebar from "../component/Sidebar";
import ChatPanel from "./ChatPanel";
import UserSidebar from "./UserSidebar";

function Messages() {
    return (
      <div className="messages__page">
          <Sidebar/>
          <ChatPanel/>
      </div>
    );
}

export default Messages;