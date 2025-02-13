import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import Link from "next/link";

// Simulated API Call (Replace with Real API)
async function getProjects() {
    // Normally, this would fetch from a database or API:
    // const res = await fetch("https://yourapi.com/projects");
    // return res.json();

    return [
        {
            id: 1,
            name: "Spring 2024 Club Applications",
            totalApplicants: 50,
            completedComparisons: 125,
            totalComparisons: 200,
            progress: 62.5
        },
        {
            id: 2, 
            name: "Fall 2023 Executive Board",
            totalApplicants: 20,
            completedComparisons: 80,
            totalComparisons: 100,
            progress: 80
        }
    ];
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
                                    <p className="text-sm text-gray-500">Total Applicants: {project.totalApplicants}</p>
                                    <p className="text-sm text-gray-500">Comparisons: {project.completedComparisons} / {project.totalComparisons}</p>
                                </div>
                                <Progress value={project.progress} className="w-full" />
                                <Button asChild className="w-full">
                                    <Link href={`/review/${project.id}`}>Continue Review</Link>
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
