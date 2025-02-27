import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import Link from "next/link";
import { Plus } from "lucide-react";

interface Project {
  id: string;
  name: string;
  totalApplicants: number;
  completedComparisons: number;
  totalComparisons: number;
}

async function getProjects(): Promise<Project[]> {
  try {
    const res = await fetch("http://localhost:8080/api/projects", {
      cache: "no-store",
    });

    if (!res.ok) {
      throw new Error("Failed to fetch projects from backend");
    }

    return res.json();
  } catch (error) {
    console.error("Error fetching projects:", error);
    return [];
  }
}

export default async function Home() {
  const projects = await getProjects();

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Active Projects</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {projects.map((project) => (
          <Card key={project.id} className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle>{project.name}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-gray-500">
                    Total Applicants: {project.totalApplicants}
                  </p>
                  <p className="text-sm text-gray-500">
                    Comparisons: {project.completedComparisons} /{" "}
                    {project.totalComparisons}
                  </p>
                </div>
                <Progress value={project.completedComparisons} className="w-full" />
                <Button asChild className="w-full">
                  <Link href={`/review/${project.id}`}>Continue Review</Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
        <Card className="flex items-center justify-center h-full hover:shadow-lg transition-shadow">
          <Button asChild variant="outline" className="w-full h-full p-8">
            <Link href="/home/add" className="flex flex-col items-center gap-4">
              <Plus className="h-12 w-12" />
              <span className="text-lg">Add New Project</span>
            </Link>
          </Button>
        </Card>
      </div>
    </div>
  );
}
