import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function HowItWorks() {
    const steps = [
        {
            title: "Upload Applicantions",
            description: "Import your applicant data into SwipeRank."
        },
        {
            title: "Compare Applicants",
            description: "Review applications in pairs, choosing your preferred candidate each time."
        },
        {
            title: "Generate Rankings",
            description: "Our algorithm processes your choices to create a comprehensive ranking."
        },
        {
            title: "Make Decisions",
            description: "Use the rankings to inform your final selection process"
        }
    ]    

    return (
        <section id="how-it-works" className="py-20 px-6 md:px-10 bg-gray-50">
            <h2 className="text-3xl md:text-4xl font-bold text-center mb-12">How It Works</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                {steps.map((step, index) => (
                    <Card key={index}>
                        <CardHeader>
                            <CardTitle className="flex items-center">
                                <span className="bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center mr-2">
                                    {index + 1}
                                </span>
                                {step.title}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <CardDescription>{step.description}</CardDescription>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </section>
    )
}