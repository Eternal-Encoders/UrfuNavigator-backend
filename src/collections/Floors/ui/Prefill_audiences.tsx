import React from "react";
import { useField, Label } from "payload/components/forms";
import { Gutter } from "payload/components/elements";
import { fromKeyToId } from "../../../utils";

function PrefillAudiences () {
    const {setValue: setWidth} = useField({ path: "width" });
    const {setValue: setHeight} = useField({ path: "height" });
    const {setValue: setAudiences} = useField({ path: "audiences" });
    const {setValue: setService} = useField({ path: "service" });

    function onChangeHandler(e: React.ChangeEvent<HTMLInputElement>) {
        if (e.target.files && e.target.files.length === 1) {
            const fileReader = new FileReader();
            fileReader.readAsText(e.target.files[0], "UTF-8");
            fileReader.onload = () => {
                const json = JSON.parse(fileReader.result as string);

                setWidth(json.width);
                setHeight(json.height);
                setAudiences(fromKeyToId(json.audiences));
                setService(json.service);
            }
        }
    }

    return (
        <Gutter>
            <Label htmlFor="all_input" label="Заполнить аудитории и сервис" />
            <input 
                onChange={onChangeHandler}
                type="file" 
                id="all_input" 
                accept="application/json" 
            />
        </Gutter>
    );
}

export default PrefillAudiences;