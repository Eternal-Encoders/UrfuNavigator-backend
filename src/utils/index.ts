import fs from 'fs'

interface IdWise {
    id: string
}

export function fromKeyToId(obj: {[id: string]: {[key: string]: any}}): {[fieldKey: string]: any}[] {
    return Object.keys(obj).map((key) => {
        return {
            ...obj[key]
        }
    });
};

export function fromIdToKey<Type extends IdWise>(obj: Type[]): {[id: string]: Type} {
    return Object.fromEntries(obj.map(data => [data.id, data]));
}