import { CollectionConfig } from 'payload/types';
import { IStair } from '../utils/interfaces';

const GraphPoints: CollectionConfig = {
    slug: 'graph_points',
    admin: {
        useAsTitle: "id",
    },
    access: {
        read: () => true
    },
    fields: [
        {
            name: "id",
            type: "text",
            required: true
        },
        {
            name: "x",
            type: "number",
            required: true
        },
        {
            name: "y",
            type: "number",
            required: true
        },
        {
            name: "links",
            type: "json",
            required: true
        },
        {
            name: "types",
            type: "json",
            required: true
        },
        {
            name: "names",
            type: "json",
            required: true
        },
        {
            name: "floor",
            type: "number",
            required: true
        },
        {
            name: "institute",
            type: "text",
            required: true
        },
        {
            name: "time",
            type: "json",
            required: true
        },
        {
            name: "info",
            type: "text"
        },
        {
            name: "description",
            type: "text"
        },
        {
            name: "menuId",
            type: "text"
        },
        {
            name: "isPassFree",
            type: "checkbox"
        },
        {
            name: "stairId",
            type: "text",
            hooks: {
                beforeChange: [async (args) => {
                    const payload = args.req.payload;

                    if (args.value) {
                        const sibilingsStairs = await payload.find({
                            collection: "stairs",
                            where: {
                                id: { not_equals: args.data.id },
                                "stairPoint.stairId" : { equals: args.data.stairId }
                            }
                        });
                        const stair = await payload.findByID({
                            collection: "stairs",
                            id: args.data.id
                        })
                        .catch(() => undefined);
                        const sibilingsIds = [];

                        if (sibilingsStairs.totalDocs !== 0) {
                            for (const sibiling of (sibilingsStairs.docs as unknown as IStair[])) {
                                sibilingsIds.push(sibiling.stairPoint.id);
                                const parsedLinks = sibiling.links.map((e) => e.id);
                                parsedLinks.push(args.data.id);
                                
                                await payload.update({
                                    collection: "stairs",
                                    id: sibiling.id,
                                    data: {
                                        id: sibiling.stairPoint.id,
                                        stairPoint: sibiling.stairPoint.id,
                                        institute: sibiling.stairPoint.institute,
                                        links: parsedLinks
                                    }
                                });
                            }
                        }
                        if (!stair) {
                            await payload.create({
                                collection: "stairs",
                                data: {
                                    id: args.data.id,
                                    stairPoint: args.data.id,
                                    institute: args.data.institute,
                                    links: sibilingsIds
                                }
                            });
                        } else {
                            await payload.update({
                                collection: "stairs",
                                id: args.data.id,
                                data: {
                                    id: args.data.id,
                                    stairPoint: args.data.id,
                                    institute: args.data.institute,
                                    links: sibilingsIds
                                }
                            })
                        }
                    }
                }]
            }
        },
    ],
    hooks: {
        afterDelete: [async (args) => {
            const payload = args.req.payload;

            if (args.doc.stairId) {
                const sibilingsStairs = await payload.find({
                    collection: "stairs",
                    where: {
                        id: { not_equals: args.doc.id },
                        "stairPoint.stairId" : { equals: args.doc.stairId }
                    }
                });
                if (sibilingsStairs.totalDocs !== 0) {
                    for (const sibiling of (sibilingsStairs.docs as unknown as IStair[])) {
                        const parsedLinks = sibiling.links.map((e) => e.id);
                        const indexOfId = parsedLinks.indexOf(args.doc.id);
                        if (indexOfId !== -1) parsedLinks.splice(indexOfId, 1);
                        
                        await payload.update({
                            collection: "stairs",
                            id: sibiling.id,
                            data: {
                                id: sibiling.stairPoint.id,
                                stairPoint: sibiling.stairPoint.id,
                                institute: sibiling.stairPoint.institute,
                                links: parsedLinks
                            }
                        });
                    }
                }

                await payload.delete({
                    collection: "stairs",
                    id: args.doc.id
                });
            }
        }]
    }
};

export default GraphPoints;