import { CollectionConfig } from 'payload/types';

const Insitutes: CollectionConfig = {
  slug: 'insitutes',
  admin: {
    useAsTitle: 'name',
  },
  access: {
    read: () => true
  },
  fields: [
    {
        name: "name",
        type: "text",
        required: true
    },
    {
      name: "displayableName",
      type: "text",
      required: true
    },
    {
        name: "minFloor",
        type: "number",
        required: true
    },
    {
        name: "maxFloor",
        type: "number",
        required: true
    },
    {
        name: "url",
        type: "text",
        required: true
    },
    {
      name: "latitude",
      type: "number",
      required: true
    },
    {
      name: "longitude",
      type: "number",
      required: true
    },
    {
        name: "icon",
        type: "upload",
        relationTo: "media",
        required: true
    }
  ],
};

export default Insitutes;