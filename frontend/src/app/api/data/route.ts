import { NextResponse } from "next/server";
// JUST FOR TESTING
export async function GET() {
  const data = [
    {
      id: 1,
      name: "people",

    },
     {
      id: 2,
      name: "people",
      
    },
    {
        id: 3,
        name: "people",
        
      },
      {
        id: 4,
        name: "people",
        
      },
      {
        id: 5,
        name: "people",
        
      },
  ];

  return NextResponse.json(data);
}
