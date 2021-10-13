package flag

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Usage prints a usage message documenting all defined command-line flags
// in a nice way formatted so u can use it to override the default flag.Usage func
// also added the to set the Usage and Example print so the Usage Msg will look much clear
// like :
//Usage:
//  gracefully_shutdown_sg_pod.exe --sgPodName=<value>
//Example:
//  gracefully_shutdown_sg_pod.exe --sgPodName=sg-1-host-1
//
//All optional flag's:
//  --kubeconfig string       Absolute path to the kubeconfig file deflate is ~.kube/config (default C:\Users\934165\.kube\config)
//  --nameSpace string        The service nameSpace (default storagepod)
//  --quiet                   To run in quiet mode without asking for confirmation (default false)
//  --sgPodName string        The storage pod name
func Usage(mustUseFlags, examples []string) func() {
	return func() {
		f := flag.CommandLine
		if len(mustUseFlags) > 0 {
			fmt.Println("Usage:")
			for _, useFlag := range mustUseFlags {
				fmt.Fprintf(f.Output(), "%s %s\n", os.Args[0], useFlag)
			}
		}
		if len(examples) > 0 {
			fmt.Println("Example:")
			for _, example := range examples {
				fmt.Fprintf(f.Output(), "%s %s\n", os.Args[0], example)
			}
			fmt.Println()
		}
		fmt.Fprintf(f.Output(), "All optional flag's:\n")
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		defer writer.Flush()
		flag.VisitAll(func(flag_ *flag.Flag) {
			if flag_.Usage == "" {
				return
			}
			s := fmt.Sprintf("  --%s", flag_.Name) //Two spaces before -; see next two comments.
			name, usage := flag.UnquoteUsage(flag_)
			if len(name) > 0 {
				s += " " + name
			}
			usage = strings.ReplaceAll(usage, "\n", "")
			if flag_.DefValue != "" {
				usage += fmt.Sprintf(" (default %v)", flag_.DefValue)
			}
			_, _ = fmt.Fprintf(writer, "%s\t    %s\t\n", s, usage)
		})
	}
}
