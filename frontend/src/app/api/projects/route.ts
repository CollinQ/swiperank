import { NextResponse } from "next/server";
// JUST FOR TESTING
export async function GET() {
  const projects = [
    {
      id: 1,
      name: "Spring 2024 Club Applications",
      totalApplicants: 50,
      completedComparisons: 125,
      totalComparisons: 200,
      progress: 62.5,
    },
    {
      id: 2,
      name: "Fall 2023 Executive Board",
      totalApplicants: 20,
      completedComparisons: 80,
      totalComparisons: 100,
      progress: 80,
    },
    {
      id: 3,
      name: "Fall 2023 Executive Board",
      totalApplicants: 20,
      completedComparisons: 80,
      totalComparisons: 100,
      progress: 80,
    },
    
    {
      id: 4,
      name: "Fall 2023 Executive Board",
      totalApplicants: 20,
      completedComparisons: 80,
      totalComparisons: 100,
      progress: 80,
    },
  ];

  return NextResponse.json(projects);
}
