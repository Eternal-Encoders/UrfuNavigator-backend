import { CollectionConfig } from 'payload/types';

const Stairs: CollectionConfig = {
  slug: 'stairs',
  admin: {
    useAsTitle: 'id',
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
      name: "stairPoint",
      type: "relationship",
      relationTo: "graph_points",
      required: true
    },
    {
      name: "institute",
      type: "text",
      required: true
    },
    {
      name: "links",
      type: "relationship",
      relationTo: "stairs",
      hasMany: true,
    }
  ],
};

export default Stairs;