import React from "react";
import { useField, Label } from "payload/components/forms";
import { Gutter } from "payload/components/elements";
import { fromKeyToId } from "../../../utils";

function PrefillGraph () {
    const {setValue: setGraph} = useField({ path: "graph" });

    function onChangeHandler(e: React.ChangeEvent<HTMLInputElement>) {
        if (e.target.files && e.target.files.length === 1) {
            const fileReader = new FileReader();
            fileReader.readAsText(e.target.files[0], "UTF-8");
            fileReader.onload = () => {
                const json = JSON.parse(fileReader.result as string);
                const graph = fromKeyToId(json.graph);
                setGraph(graph);
            }
        }
    }

    return (
        <Gutter>
            <Label htmlFor="all_input" label="Заполнить граф" />
            <input 
                onChange={onChangeHandler}
                type="file" 
                id="all_input" 
                accept="application/json" 
            />
        </Gutter>
    );
}

export default PrefillGraph;