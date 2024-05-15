import React from "react";
import payload from "payload";
import "./devOrTest.css";

function DevOrTestText () {
    console.log(process.env.PORT_ENV);

    return (
        <div className="devortest-container">
            <h3 className="devortest-text">
                {window.location.port === "8443" ? "ПРОДАКШЕН": "ТЕСТОВЫЙ"}
            </h3>
        </div>
    )
}

export default DevOrTestText;