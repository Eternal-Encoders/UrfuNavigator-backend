import { CollectionConfig } from 'payload/types';
import PrefillAll from './ui/Prefill_all';
import PrefillAudiences from './ui/Prefill_audiences';
import PrefillGraph from './ui/Prefill_graph';
import payload from 'payload';
import Download from './ui/Download';

const Floors: CollectionConfig = {
  slug: 'floors',
  admin: {
    useAsTitle: "institute",
    defaultColumns: [
        "id",
        "institute",
        "floor"
    ]
  },
  access: {
    read: () => true
  },
  fields: [
    {
        label: "Скачать граф",
        name: "download",
        type: "ui",
        admin: {
            components: {
                Field: Download,
            },
            condition: (data) => {
                return data.createdAt
            }
        }
    },
    {
        label: "Загрузить из файла",
        type: "collapsible",
        fields: [
            {
                name: "prefill_all",
                type: "ui",
                admin: {
                    components: {
                        Field: PrefillAll,
                    }
                }
            },
            {
                name: "prefill_draw",
                type: "ui",
                admin: {
                    components: {
                        Field: PrefillAudiences
                    }
                }
            },
            {
                name: "prefill_graph",
                type: "ui",
                admin: {
                    components: {
                        Field: PrefillGraph
                    }
                }
            },
        ]
    },
    {
        name: "institute",
        type: "text",
        required: true
    },
    {
        name: "floor",
        type: "number",
        required: true
    },
    {
        name: "width",
        type: "number",
        required: true
    },
    {
        name: "height",
        type: "number",
        required: true
    },
    {
        name: "audiences",
        type: "json",
        required: true,
        admin: {
            hidden: true
        }
    },
    {
        name: "service",
        type: "json",
        required: true,
        admin: {
            hidden: true
        }
    },
    {
        name: "graph",
        type: "relationship",
        relationTo: "graph_points",
        hasMany: true,
        required: true,
        admin: {
            hidden: true
        },
        hooks: {
            beforeChange: [async (args) => {
                if(args.value.length !== 0 && typeof args.value[0] === "string") {
                    return args.value
                }

                const payload = args.req.payload;
                
                if (args.originalDoc.graph) {
                    for (const prevPoint of args.originalDoc.graph) {
                        await payload.delete({
                            collection: "graph_points",
                            id: prevPoint
                        });
                    }
                }

                for (const point of args.value) {
                    await payload.create({
                        collection: "graph_points",
                        data: {
                            ...point
                        }
                    });
                }
                return args.value.map(({id}) => {
                    return id
                });
                
            }]
        }
    }
  ],
  hooks: {
    afterDelete: [async (args) => {
        const payload = args.req.payload;
        for (const point of args.doc.graph) {
            await payload.delete({
                collection: "graph_points",
                id: point.id
            });
        }
    }]
  }
};

export default Floors;