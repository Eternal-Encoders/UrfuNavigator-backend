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

export function loadEnvOrFile(name: string): string {
    let data = process.env[name]
    if (!data) {
        const path = process.env[`${name}_FILE`]
        if (!path) {
            throw Error(`No ${name} specified`)
        }
        
        try {
            data = fs.readFileSync(path, 'utf8');
        } catch (err) {
            throw err
        }
    }

    return data
}