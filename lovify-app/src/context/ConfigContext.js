import {createContext, useContext} from "react";

const ConfigContext = createContext();

export const ConfigProvider = ({children}) => {
    const config = {
        apiUrl: process.env.REACT_APP_API_URL || "",
    };

    return (
        <ConfigContext.Provider value={config}>
            {children}
        </ConfigContext.Provider>
    );
};

export const useConfig = () => {
    const context = useContext(ConfigContext)
    if (!context) {
        throw new Error('useConfig debe usarse dentro de un ConfigProvider');
    }
    return context
}