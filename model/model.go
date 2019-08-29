package model

// SequenceAnnotation represents a single row describing a
// mutation of a specific nucleotide sequence. Here's an
// example of what a row typically looks like from a version 0.27.1
// of the breseq application.
//
// See for more details: http://barricklab.org/twiki/pub/Lab/ToolsBacterialGenomeResequencing/documentation/output.html
//
// seq_id    | position | mutation | freq | annotation            | gene     | description
// NC_012345 | 12,345   | +G       | 100% | intergenic (-123/+12) | ABC01234 | lipoprotein
// NC_012345 | 65,431   | +A       | 6.0% | V12A (GTG&rarr;GGG)   | ABC05678 | hypothetical protein
// etc.
type SequenceAnnotation struct {
	// UniqueId is an application generated string that uniquely
	// identifies this sequence annotation.
	UniqueId string
	// SequenceId is the identifier for the reference sequence
	// with the mutation.
	SequenceId string
	// Position in the reference sequence of the mutation.
	Position string
	// Generation serves two functions:
	// (1) Groups all sequence annotations by generation.
	// (2) Acts as a timestamp to differentiate the sequence annotations
	//     with the same sequence id and position.
	Generation string
	// Application is the name of the application this
	// annotation came from.
	Application string
	// AppVersion is the version of the application this
	// annotation came from.
	AppVersion string
	// Mutation is a description, usually of how nucleotides
	// are added, substituted, or deleted.
	Mutation string
	// Frequency is a percentage field of how often this
	// mutation occurs.
	Frequency string
	// Annotation is a more detailed description of the mutation.
	Annotation string
	// Gene is a space-delimited list of genes affected by the mutation.
	Gene string
	// Description is a qualitative description of the genes affected.
	Description string
}
