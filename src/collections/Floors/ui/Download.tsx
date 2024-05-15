import React from "react";
import { useField, Label } from "payload/components/forms";
import { Gutter } from "payload/components/elements";
import { fromIdToKey } from "../../../utils";
import { IMapDownload, IService, IAuditorium, IGraphPoint } from "../../../utils/interfaces";
import payload from "payload";

function Download () {
    const {value: institute} = useField<string>({ path: "institute" });
    const {value: floor} = useField<number>({ path: "floor" });
    const {value: width} = useField<number>({ path: "width" });
    const {value: height} = useField<number>({ path: "height" });
    const {value: audiences} = useField<IAuditorium[]>({ path: "audiences" });
    const {value: service} = useField<IService[]>({ path: "service" });
    const {value: graphIds} = useField<string[]>({ path: "graph" });

    function onClickHandler(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
        e.preventDefault();
        e.stopPropagation();

        const link = document.createElement("a");
        const promises: Promise<IGraphPoint>[] = []

        graphIds.forEach(e => {
            promises.push(fetch(`/api/graph_points/${e}`)
                .then(res => res.json())
                .then((data) => {
                    data.updatedAt = undefined;
                    data.createdAt = undefined;
                    return data
                }));
        });

        Promise.all(promises)
            .then((graph) => {
                const obj: IMapDownload = {
                    service,
                    audiences: fromIdToKey<IAuditorium>(audiences),
                    graph: fromIdToKey<IGraphPoint>(graph),
                    institute,
                    floor,
                    width,
                    height
                }
        
                const newFile = new Blob([JSON.stringify(obj)], { type: "application/json" });
                link.href = URL.createObjectURL(newFile);
                link.download = `${institute}_${floor}.json`;
                link.click();
                link.remove();
            })
    }

    return (
        <Gutter>
            <Label />
            <button onClick={onClickHandler}>Скачать граф</button>
        </Gutter> 
    )
}

export default Download;